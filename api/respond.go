package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func encodeBody(writer http.ResponseWriter, request *http.Request, v interface{}) error {
	return json.NewEncoder(writer).Encode(v)
}

func respondErr(writer http.ResponseWriter, request *http.Request, statusCode int, args ...interface{}) {
	respond(writer, request, statusCode, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}

func respondHTTPErr(w http.ResponseWriter, r *http.Request, status int) {
	respondErr(w, r, status, http.StatusText(status))
}

func respond(writer http.ResponseWriter, request *http.Request, statusCode int, data interface{}) {
	writer.WriteHeader(statusCode)
	if data != nil {
		encodeBody(writer, request, data)
	}
}
