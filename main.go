package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := getenv("MAX_REQUESTS", 500)
	t := getenv("TOKEN", "token")
	r := mux.NewRouter()
	r.HandleFunc("/", handler(l)).
		Methods("POST")
	r.HandleFunc("/healthz", httpHealth(l)).
		Methods("GET")
	r.HandleFunc("/reset", reset(t)).
		Methods("PUT")
	logger := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8080", logger))
}
