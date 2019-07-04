package token

import (
	"encoding/json"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/goutil/assert"
	"testing"
)

type ctx struct {
	Name string `json:"name"`
}

func (c *ctx) ToMap() goutil.Map {
	return goutil.Struct2Map(c)
}

func (c *ctx) LoadFromMap(data map[string]interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, &c)
}

func TestParse(t *testing.T) {
	Init(&Config{Secret: "secret"})

	c := &ctx{Name: "tom"}
	token, err := GenerateToken(c)
	assert.NoError(t, err)

	p := &ctx{}
	err = Parse(token, conf.Secret, p)
	assert.NoError(t, err)

	assert.Equal(t, c.Name, p.Name)
}
