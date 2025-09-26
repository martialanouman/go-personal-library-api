package helpers

import (
	"encoding/json"
	"net/http"
)

type Envelop map[string]any

func WriteJson(w http.ResponseWriter, status int, data Envelop) error {
	js, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
