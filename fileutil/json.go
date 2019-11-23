package fileutil

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func WriteJSONFile(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(data)
}

func ReadJSONFile(filename string, out interface{}) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, out)
}