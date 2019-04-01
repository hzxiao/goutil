package version

import (
	"github.com/hzxiao/goutil/assert"
	"testing"
)

func TestPrint(t *testing.T) {
	err := Print()
	assert.NoError(t, err)
}
