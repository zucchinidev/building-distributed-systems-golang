package main

import (
	"context"
	"github.com/joeshaw/envdecode"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

var config struct {
	MongoDBName         string `env:"BDSG_MONGO_DB_NAME,required"`
	MongoCollectionName string `env:"BDSG_MONGO_COLLECTION_NAME,required"`
	Addr                string `env:"BDSG_API_ADDR,required"`
	MongoService        string `env:"BDSG_MONGO_SERVICE,required"`
}

func init() {
	if err := envdecode.Decode(&config); err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.Println("Dialign mongodb", config.MongoService)
	db, err := mgo.Dial(config.MongoService)
	if err != nil {
		log.Fatalln("failed to connect to mongodb: ", err)
	}
	defer db.Close()

	server := &Server{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withAPIKey(server.handlePolls)))
	log.Println("Starting web server on ", config.Addr)
	http.ListenAndServe(config.Addr, mux)
	log.Println("Stopping...")
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
