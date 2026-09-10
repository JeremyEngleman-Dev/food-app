package utils

import (
	"encoding/json"
	"net/http"
)

func HttpJsonResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)

	_, err = w.Write(body)

	return err
}
