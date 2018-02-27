package goutil

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestMap_Get(t *testing.T) {
	m := Map{}

	m.Set("string", "1")
	stringInterf := m.Get("string")
	_string, ok := stringInterf.(string)
	assert.True(t, ok)
	assert.Equal(t, "1", _string)

	m.Set("true", true)
	m.Set("false", false)
	assert.True(t, m.GetBool("true"))
	assert.False(t, m.GetBool("false"))

	m.Set("int", int(123))
	assert.Equal(t, int64(123), m.GetInt64("int"))

	m.Set("int8", int8(12))
	assert.Equal(t, int64(12), m.GetInt64("int8"))

	m.Set("int16", int16(123))
	assert.Equal(t, int64(123), m.GetInt64("int16"))

	m.Set("int32", int32(123))
	assert.Equal(t, int64(123), m.GetInt64("int32"))

	m.Set("int64", int64(123))
	assert.Equal(t, int64(123), m.GetInt64("int64"))

	m.Set("uint", uint(123))
	assert.Equal(t, int64(123), m.GetInt64("uint"))

	m.Set("uint8", uint8(123))
	assert.Equal(t, int64(123), m.GetInt64("uint8"))

	m.Set("uint16", uint16(123))
	assert.Equal(t, int64(123), m.GetInt64("uint16"))

	m.Set("uint32", uint32(123))
	assert.Equal(t, int64(123), m.GetInt64("uint32"))

	m.Set("uint64", uint64(123))
	assert.Equal(t, int64(123), m.GetInt64("uint64"))

	m.Set("float32", float32(123))
	assert.Equal(t, float64(123), m.GetFloat64("float32"))

	m.Set("float64", float64(123))
	assert.Equal(t, float64(123), m.GetFloat64("float64"))

	m.Set("map", Map{"int": 123})
	assert.Len(t, m.GetMap("map"), 1)

	m.Set("int64Array", []int64{1, 2, 3})
	int64Array := m.GetInt64Array("int64Array")
	assert.Len(t, int64Array, 3)

	m.Set("stringArray", []string{"1", "2", "3"})
	assert.Len(t, m.GetStringArray("stringArray"), 3)

	m.Set("float64Array", []float64{1, 2, 3})
	assert.Len(t, m.GetFloat64Array("float64Array"), 3)
	m.Set("mapArray", []Map{
		{"int": 123},
		{"int": 123},
		{"int": 123},
	})
	assert.Len(t, m.GetMapArray("mapArray"), 3)

	assert.Len(t, m, 20)
}

func TestMap_GetP(t *testing.T) {
	m := Map{}
	m.Set("map", Map{
		"int":         123,
		"true":        true,
		"float64":     float64(123),
		"string":      "123",
		"stringArray": []string{"1", "2"},
		"map":         Map{"int": 123},
		"mapArray": []Map{
			Map{"int": 123},
			Map{"int": 123},
		},
		"intArray":     []int64{1, 2, 3},
		"float64Array": []float64{1, 2},
	})

	assert.Equal(t, int64(123), m.GetInt64P("map/int"))
	assert.True(t, m.GetBoolP("map/true"))
	assert.Equal(t, float64(123), m.GetFloat64P("map/float64"))
	assert.Equal(t, "123", m.GetStringP("map/string"))
	assert.Len(t, m.GetStringArrayP("map/stringArray"), 2)
	assert.Equal(t, "2", m.GetStringP("map/stringArray/1"))
	assert.Len(t, m.GetMapP("map/map"), 1)
	assert.Equal(t, int64(123), m.GetInt64P("map/map/int"))
	assert.Len(t, m.GetMapArrayP("map/mapArray"), 2)
	assert.Equal(t, int64(123), m.GetInt64P("map/mapArray/0/int"))
	assert.Len(t, m.GetInt64ArrayP("map/intArray"), 3)
	assert.Len(t, m.GetFloat64ArrayP("map/float64Array"), 2)
}
