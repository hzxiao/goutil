package httputil

import (
	"encoding/json"
	"encoding/xml"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

type Person struct {
	XMLName   xml.Name `xml:"person"`
	Id        int      `xml:"id,attr"`
	FirstName string   `xml:"name>first"`
	LastName  string   `xml:"name>last"`
	Age       int      `xml:"age"`
	Height    float32  `xml:"height,omitempty"`
	Married   bool
	Comment   string `xml:",comment"`
}

var person = `<person id="13"><name><first>John</first><last>Doe</last></name><age>42</age><Married>false</Married></person>`

func TestFormatValue(t *testing.T) {
	var by []byte
	err := formatValue(ReturnBytes, []byte{1, 2, 3}, &by)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, by)

	var by2 []byte
	err = formatValue(ReturnBytes, nil, &by2)
	assert.NoError(t, err)
	assert.Nil(t, by2)

	var str string
	err = formatValue(ReturnString, []byte("hello"), &str)
	assert.NoError(t, err)
	assert.Equal(t, "hello", str)

	var str2 string
	err = formatValue(ReturnString, nil, &str2)
	assert.NoError(t, err)
	assert.Equal(t, "", str2)

	var m goutil.Map
	err = formatValue(ReturnJSON, []byte(`{"key":"value"}`), &m)
	assert.NoError(t, err)
	assert.Equal(t, goutil.Map{"key": "value"}, m)

	var p Person
	err = formatValue(ReturnXML, []byte(person), &p)
	assert.NoError(t, err)
	assert.Equal(t, Person{XMLName: xml.Name{Local: "person"}, Id: 13, FirstName: "John", LastName: "Doe", Age: 42}, p)
}

func unformatValue(returnType int, v interface{}) ([]byte, error) {
	switch returnType {
	case ReturnBytes:
		b, _ := v.([]byte)
		return b, nil
	case ReturnString:
		s, _ := v.(string)
		return []byte(s), nil
	case ReturnJSON:
		return json.Marshal(v)
	case ReturnXML:
		return xml.Marshal(v)
	}
	return nil, nil
}

func TestClient_Get(t *testing.T) {
	p := Person{XMLName: xml.Name{Local: "person"}, Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	var tables = []struct {
		ReturnErr   bool
		ReturnType  int
		ReturnValue interface{}
		ReturnCode  int
		CheckResult interface{}
	}{
		{false, ReturnBytes, []byte{1, 2, 3}, 200, nil},
		{false, ReturnString, "str", 200, nil},
		{false, ReturnJSON, map[string]interface{}{"k": "v"}, 200, nil},
		{false, ReturnXML, p, 200, nil},
		{true, ReturnBytes, nil, 200, nil},
		{false, ReturnBytes, nil, 500, nil},
	}

	var srv *httptest.Server
	var handler = func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)
		index, err := strconv.Atoi(r.FormValue("index"))
		assert.NoError(t, err)
		c := tables[index]
		if c.ReturnErr {
			srv.CloseClientConnections()
			return
		}
		if c.ReturnCode != 200 {
			w.WriteHeader(c.ReturnCode)
			return
		}

		res, err := unformatValue(c.ReturnType, c.ReturnValue)
		assert.NoError(t, err)
		if c.ReturnType == ReturnJSON {
			w.Header().Set("Content-Type", "application/json")
		}
		if c.ReturnType == ReturnXML {
			w.Header().Set("Content-Type", "application/xml")
		}

		w.Write(res)
	}
	srv = httptest.NewServer(http.HandlerFunc(handler))

	for i := range tables {
		if tables[i].ReturnType == ReturnXML {
			tables[i].CheckResult = &Person{}
		}
		err := HTTPClient.Get(srv.URL+"/?index="+strconv.Itoa(i), tables[i].ReturnType, &tables[i].CheckResult)
		if tables[i].ReturnErr || tables[i].ReturnCode != 200 {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		if tables[i].ReturnType == ReturnXML {
			tables[i].CheckResult = *(tables[i].CheckResult.(*Person))
		}
		assert.Equal(t, tables[i].ReturnValue, tables[i].CheckResult)
	}
}

func TestClient_PostForm(t *testing.T) {
	p := Person{XMLName: xml.Name{Local: "person"}, Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	var tables = []struct {
		ReturnErr   bool
		ReturnType  int
		ReturnValue interface{}
		ReturnCode  int
		CheckResult interface{}
	}{
		{
			ReturnErr:   false,
			ReturnType:  ReturnBytes,
			ReturnValue: []byte{1, 2, 3},
			ReturnCode:  200,
		},
		{
			ReturnErr:   false,
			ReturnType:  ReturnString,
			ReturnValue: "hello",
			ReturnCode:  200,
		},
		{
			ReturnErr:   false,
			ReturnType:  ReturnJSON,
			ReturnValue: map[string]interface{}{"k": "v"},
			ReturnCode:  200,
		},
		{
			ReturnErr:   false,
			ReturnType:  ReturnXML,
			ReturnValue: p,
			ReturnCode:  200,
		},
		{
			ReturnErr: true,
		},
		{
			ReturnCode: 500,
		},
	}

	var srv *httptest.Server
	var handler = func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		assert.NoError(t, err)
		index, err := strconv.Atoi(r.PostFormValue("index"))
		assert.NoError(t, err)
		c := tables[index]
		if c.ReturnErr {
			srv.CloseClientConnections()
			return
		}
		if c.ReturnCode != 200 {
			w.WriteHeader(c.ReturnCode)
			return
		}

		res, err := unformatValue(c.ReturnType, c.ReturnValue)
		assert.NoError(t, err)
		if c.ReturnType == ReturnJSON {
			w.Header().Set("Content-Type", "application/json")
		}
		if c.ReturnType == ReturnXML {
			w.Header().Set("Content-Type", "application/xml")
		}
		w.Write(res)
	}
	srv = httptest.NewServer(http.HandlerFunc(handler))

	for i := range tables {
		if tables[i].ReturnType == ReturnXML {
			tables[i].CheckResult = &Person{}
		}
		err := HTTPClient.PostForm(srv.URL, goutil.Map{"index": strconv.Itoa(i)}, tables[i].ReturnType, &tables[i].CheckResult)
		if tables[i].ReturnErr || tables[i].ReturnCode != 200 {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		if tables[i].ReturnType == ReturnXML {
			tables[i].CheckResult = *(tables[i].CheckResult.(*Person))
		}
		assert.Equal(t, tables[i].ReturnValue, tables[i].CheckResult)
	}
}

func TestClient_PostJSON(t *testing.T) {
	var tables = []struct {
		ReturnErr   bool
		ReturnType  int
		ReturnValue interface{}
		ReturnCode  int
		CheckResult interface{}
	}{
		{
			ReturnErr:   false,
			ReturnType:  ReturnBytes,
			ReturnValue: []byte{1, 2, 3},
			ReturnCode:  200,
		},
		{
			ReturnErr:   false,
			ReturnType:  ReturnString,
			ReturnValue: "hello",
			ReturnCode:  200,
		},
		{
			ReturnErr:   false,
			ReturnType:  ReturnJSON,
			ReturnValue: map[string]interface{}{"k": "v"},
			ReturnCode:  200,
		},
		{
			ReturnErr: true,
		},
		{
			ReturnCode: 500,
		},
	}

	var srv *httptest.Server
	var handler = func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()

		var data goutil.Map
		err = json.Unmarshal(buf, &data)
		assert.NoError(t, err)
		c := tables[int(data.GetInt64("index"))]
		if c.ReturnErr {
			srv.CloseClientConnections()
			return
		}
		if c.ReturnCode != 200 {
			w.WriteHeader(c.ReturnCode)
			return
		}

		res, err := unformatValue(c.ReturnType, c.ReturnValue)
		assert.NoError(t, err)
		if c.ReturnType == ReturnJSON {
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write(res)
	}
	srv = httptest.NewServer(http.HandlerFunc(handler))

	for i := range tables {
		err := HTTPClient.PostJSON(srv.URL, goutil.Map{"index": strconv.Itoa(i)}, tables[i].ReturnType, &tables[i].CheckResult)
		if tables[i].ReturnErr || tables[i].ReturnCode != 200 {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		assert.Equal(t, tables[i].ReturnValue, tables[i].CheckResult)
	}
}

type xmlReqData struct {
	XMLName   xml.Name `xml:"person"`
	Index        int      `xml:"index"`
}

func TestClient_PostXML(t *testing.T) {
	p := Person{XMLName: xml.Name{Local: "person"}, Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	var tables = []struct {
		ReturnErr   bool
		ReturnType  int
		ReturnValue interface{}
		ReturnCode  int
		CheckResult interface{}
	}{
		{
			ReturnErr:   false,
			ReturnType:  ReturnBytes,
			ReturnValue: []byte{1, 2, 3},
			ReturnCode:  200,
		},
		{
			ReturnErr:   false,
			ReturnType:  ReturnString,
			ReturnValue: "hello",
			ReturnCode:  200,
		},
		{
			ReturnErr:   false,
			ReturnType:  ReturnXML,
			ReturnValue: p,
			ReturnCode:  200,
		},
		{
			ReturnErr: true,
		},
		{
			ReturnCode: 500,
		},
	}

	var srv *httptest.Server
	var handler = func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()

		var data xmlReqData
		err = xml.Unmarshal(buf, &data)
		assert.NoError(t, err)
		c := tables[data.Index]
		if c.ReturnErr {
			srv.CloseClientConnections()
			return
		}
		if c.ReturnCode != 200 {
			w.WriteHeader(c.ReturnCode)
			return
		}

		res, err := unformatValue(c.ReturnType, c.ReturnValue)
		assert.NoError(t, err)
		if c.ReturnType == ReturnJSON {
			w.Header().Set("Content-Type", "application/json")
		}
		if c.ReturnType == ReturnXML {
			w.Header().Set("Content-Type", "application/xml")
		}
		w.Write(res)
	}
	srv = httptest.NewServer(http.HandlerFunc(handler))

	for i := range tables {
		if tables[i].ReturnType == ReturnXML {
			tables[i].CheckResult = &Person{}
		}
		reqData := &xmlReqData{Index: i}
		err := HTTPClient.PostXML(srv.URL, reqData, tables[i].ReturnType, &tables[i].CheckResult)
		if tables[i].ReturnErr || tables[i].ReturnCode != 200 {
			assert.Error(t, err)
			continue
		}
		assert.NoError(t, err)
		if tables[i].ReturnType == ReturnXML {
			tables[i].CheckResult = *(tables[i].CheckResult.(*Person))
		}
		assert.Equal(t, tables[i].ReturnValue, tables[i].CheckResult)
	}
}

func TestClient_Download(t *testing.T) {
	var srv *httptest.Server

	srv = httptest.NewServer(http.FileServer(http.Dir("../testdata/")))

	l, err := HTTPClient.Download(srv.URL+"/test.txt","./test.txt")
	assert.NoError(t, err)

	assert.NotEqual(t, 0, int(l))

	err = os.Remove("./test.txt")
	assert.NoError(t, err)
}
