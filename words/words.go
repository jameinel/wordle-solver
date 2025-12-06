package words

import (
	"encoding/json"
	"io"
	"os"
)

type Words struct {
	Dictionary []string `json:"Ta"`
	Solutions  []string `json:"La"`
}

func readFromJSON(r io.Reader) (*Words, error) {
	var w Words
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

func ReadWordsFile(filename string) (*Words, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return readFromJSON(f)
}
