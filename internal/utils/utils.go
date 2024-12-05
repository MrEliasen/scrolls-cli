package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func Unmarshal[T any](r *http.Response) (T, error) {
	t := new(T)

	p, err := io.ReadAll(r.Body)
	if err != nil {
		return *t, err
	}

	err = json.Unmarshal(p, &t)

	return *t, nil
}
