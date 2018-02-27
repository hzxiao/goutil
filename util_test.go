package goutil

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
	"time"
)

type TestStruct struct {
	Name string
	Age  int64
}

func TestStruct2Map(t *testing.T) {
	//1
	t1 := TestStruct{Name: "name", Age: 2}
	m1 := Struct2Map(t1)
	if m1.GetString("Name") != t1.Name || m1.GetInt64("Age") != t1.Age {
		t.Error("error")
	}

	//2 struct pointer
	t2 := &t1
	m2 := Struct2Map(t2)
	if m2.GetString("Name") != t2.Name || m2.GetInt64("Age") != t2.Age {
		t.Error("error")
	}

	//3 string
	m3 := Struct2Map("string")
	if m3 != nil {
		t.Error("error")
	}
}

func BenchmarkStruct2Map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Struct2Map(TestStruct{Name: "123", Age: 20})
	}
}

func TestStruct2Json(t *testing.T) {
	//1
	t1 := &TestStruct{Name: "t1", Age: 2}
	json1 := Struct2Json(t1)
	if json1 == "" || json1 == "null" {
		t.Error("struct point -> json error")
	}
	//2
	t2 := TestStruct{Name: "t2", Age: 2}
	json2 := Struct2Json(t2)
	if json2 == "" || json2 == "null" {
		t.Error("struct -> json error")
	}
}

func TestBoolE(t *testing.T) {
	v1 := true
	b1, err := BoolE(v1)
	assert.NoError(t, err)
	assert.True(t, b1)

	v2 := false
	b2, err := BoolE(v2)
	assert.NoError(t, err)
	assert.False(t, b2)

	//test error
	_, err = BoolE(nil)
	assert.Error(t, err)

	_, err = BoolE(123)
	assert.Error(t, err)
}

func TestBool(t *testing.T) {
	assert.True(t, Bool(true))
	assert.False(t, Bool(false))
	assert.False(t, Bool(123))
}

func TestInt64E(t *testing.T) {
	//time
	_time, err := Int64E(time.Now())
	assert.NoError(t, err)
	assert.NotEqual(t, 0, _time)

	//int
	_int, err := Int64E(int(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _int)

	//int8
	_int8, err := Int64E(int8(12))
	assert.NoError(t, err)
	assert.Equal(t, int64(12), _int8)

	//int16
	_int16, err := Int64E(int16(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _int16)

	//int32
	_int32, err := Int64E(int32(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _int32)

	//int64
	_int64, err := Int64E(int64(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _int64)

	//uint
	_uint, err := Int64E(uint(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _uint)

	//uint8
	_uint8, err := Int64E(uint8(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _uint8)

	//uint16
	_uint16, err := Int64E(uint16(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _uint16)

	//uint32
	_uint32, err := Int64E(uint32(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _uint32)

	//uint64
	_uint64, err := Int64E(uint64(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _uint64)

	//float32
	_float32, err := Int64E(float32(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _float32)

	//float64
	_float64, err := Int64E(float64(123))
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _float64)

	//string
	_string, err := Int64E("123")
	assert.NoError(t, err)
	assert.Equal(t, int64(123), _string)

	//test error
	_, err = Int64E(nil)
	assert.Error(t, err)

	_, err = Int64E(true)
	assert.Error(t, err)

	_, err = Int64E(TestStruct{})
	assert.Error(t, err)

	_, err = Int64E("aa")
	assert.Error(t, err)
}
