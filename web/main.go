package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	webServerAddr := os.Getenv("BDSG_WEB_SERVER_ADDR")
	if webServerAddr == "" {
		log.Fatal("BDSG_WEB_SERVER_ADDR environment variable is mandatory")
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))
	log.Println("Serving website at: ", webServerAddr)
	http.ListenAndServe(webServerAddr, mux)
}
