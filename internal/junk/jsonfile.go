package junk

import (
	"encoding/json"
	"os"
)

func ParseJSONFile(fileName string, out interface{}) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, out)
}
