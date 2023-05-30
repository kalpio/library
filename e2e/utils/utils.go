package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func GetBodyBytes(body io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func MustMarshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return data, nil
}
