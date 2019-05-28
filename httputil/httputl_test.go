package httputil

import (
	"encoding/xml"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"net/http"
	"net/http/httptest"
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
	Comment string `xml:",comment"`
}

var person = `<person id="13"><name><first>John</first><last>Doe</last></name><age>42</age><Married>false</Married></person>`

func TestFormatValue(t *testing.T) {
	var by []byte
	err := formatValue(ReturnBytes, []byte{1, 2,3}, &by)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2,3}, by)

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
	assert.Equal(t, goutil.Map{"key":"value"}, m)

	var p Person
	err = formatValue(ReturnXML, []byte(person), &p)
	assert.NoError(t, err)
	assert.Equal(t, Person{XMLName: xml.Name{Local:"person"}, Id: 13, FirstName: "John", LastName: "Doe", Age: 42}, p)
}

func handler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		w.Write([]byte("hello"))
		return
	}
}

func TestClient_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(handler))

	var str string
	err := HTTPClient.Get(srv.URL, ReturnString, &str)
	assert.NoError(t,err)
	assert.Equal(t, "hello", str)
}