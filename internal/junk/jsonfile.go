package junk

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileContentsError struct {
	FileName   string
	Underlying error
}

func (f *FileContentsError) Error() string {
	return fmt.Sprintf("failed to use %q because %s", f.FileName, f.Underlying.Error())
}

func ParseJSONFile(fileName string, out interface{}) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return &FileContentsError{
			FileName:   fileName,
			Underlying: err,
		}
	}

	return json.Unmarshal(bytes, out)
}
