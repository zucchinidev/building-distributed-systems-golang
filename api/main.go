package main

import (
	"context"
	"gopkg.in/mgo.v2"
	"net/http"
)

func main() {

}

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

func APIKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(writer, request, http.StatusUnauthorized, "invalid API Key")
			return
		}

		ctx := context.WithValue(request.Context(), contextKeyAPIKey, key)
		fn(writer, request.WithContext(ctx))
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(writer, request)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

type Server struct {
	db *mgo.Session
}
