package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func UnmarshalResp[T any](r *http.Response) (T, error) {
	p, err := io.ReadAll(r.Body)
	if err != nil {
		return *new(T), err
	}

	return Unmarshal[T](p)
}

func Unmarshal[T any](p []byte) (T, error) {
	var t T
	err := json.Unmarshal(p, &t)
	return t, err
}

func Marshal(data any) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(data)
	return buf.Bytes(), err
}
