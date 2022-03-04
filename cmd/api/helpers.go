package main

import (
	"encoding/json"
	"net/http"
)

func (application *application) writeJson(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	response = append(response, '\n')

	// if headers have been provided, set them
	for key, value := range headers {
		w.Header().Set(key, value[0])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write((response))

	return nil
}
