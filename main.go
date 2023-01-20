package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	l int
	C counter
)

type counter struct {
	mu sync.Mutex
	n  int
}

func (c *counter) Add() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *counter) Get() int {
	c.mu.Lock()
	n := c.n
	c.mu.Unlock()
	return n
}

func (c *counter) Reset() {
	c.mu.Lock()
	c.n = 0
	c.mu.Unlock()
}

func getenv(k string, d int) int {
	v := os.Getenv(k)
	if len(v) == 0 {
		return d
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Invalid Value, %s not a valid integer: %v", k, err)
	}
	return i
}

func serve() bool {
	return C.Get() < l
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	jm, err := json.Marshal(r.PostForm)
	if err != nil || !serve() {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "why did you do that?")
		C.Add()
		return
	}
	fmt.Fprintf(w, "%v", string(jm))
	C.Add()
}

func httpHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !serve() {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"FAIL","requests":"%v"}`, C.Get())
			return
		}
		fmt.Fprintf(w, `{"status":"ok","requests":"%v"}`, C.Get())
	}
}

func main() {
	l = getenv("MAX_REQUESTS", 500)
	r := mux.NewRouter()
	r.HandleFunc("/", handler).
		Methods("POST")
	r.HandleFunc("/healthz", httpHealth()).
		Methods("GET")
	logger := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(":8080", logger))
}
