package utils

import (
	"bytes"
	"io"
)

func GetBodyBytes(body io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
