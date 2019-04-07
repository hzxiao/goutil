package assert

import (
	"errors"
	"testing"
)

func TestEqual(t *testing.T) {
	Equal(t, "Hello", "Hello")

	Equal(t, 123, 123)

	Equal(t, 123.5, 123.5)

	Equal(t, []byte("Hello World"), []byte("Hello World"))

	Equal(t, nil, nil)

	Equal(t, int32(123), int32(123))

	Equal(t, uint64(123), uint64(123))

	Equal(t, &struct{}{}, &struct{}{})
}

func TestNil(t *testing.T) {
	Nil(t, nil)

	Nil(t, (*struct{})(nil))
}

func TestTrue(t *testing.T) {
	True(t, true)
}

func TestFalse(t *testing.T) {
	False(t, false)
}

func TestNoError(t *testing.T) {
	var err error

	True(t, NoError(t, err))
}

func TestError(t *testing.T) {
	var err error

	err = errors.New("some err")
	True(t, Error(t, err))
}
