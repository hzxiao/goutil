package container

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestSafeSlice_Append(t *testing.T) {
	s := NewSafeSlice(0)
	for i := 0; i < 1000000; i++ {
		s.Append(i)
	}

	assert.Equal(t, 1000000, s.Len())
}

