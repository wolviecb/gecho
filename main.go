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

// Global counter
var C counter

// Counter mutex for counter
type counter struct {
	mu sync.Mutex
	n  int
}

// Add value to counter
func (c *counter) Add() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

// Get value from counter
func (c *counter) Get() int {
	c.mu.Lock()
	n := c.n
	c.mu.Unlock()
	return n
}

// Reset counter
func (c *counter) Reset() {
	c.mu.Lock()
	c.n = 0
	c.mu.Unlock()
}

// getenv reads an environment  variable named k and returns it as type D
func getenv[D ~string | int](k string, d D) D {
	v := os.Getenv(k)
	if len(v) == 0 {
		return d
	}
	var r any
	switch any(d).(type) {
	case string:
		r = v
	case int:
		i, err := strconv.Atoi(v)
		if err != nil {
		}
		r = i
	default:
		log.Fatalf("Invalid Value, %s not a valid", k)
	}
	return r.(D)
}

// serve evaluates if the limit of requests is reached
func serve(l int) bool {
	return C.Get() < l
}

// handler generates the echo server response
func handler(l int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		jm, err := json.Marshal(r.PostForm)
		if err != nil || !serve(l) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "why did you do that?")
			C.Add()
			return
		}
		fmt.Fprintf(w, "%v", string(jm))
		C.Add()
	}
}

// httpHealth is a bad implementation of a health check
func httpHealth(l int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !serve(l) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"FAIL","requests":"%v"}`, C.Get())
			return
		}
		fmt.Fprintf(w, `{"status":"ok","requests":"%v"}`, C.Get())
	}
}

// reset resets the request counter
func reset(rt string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		t := r.FormValue("TOKEN")
		if len(t) == 0 || t != rt {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad request, invalid token")
			return
		}
		C.Reset()
		fmt.Fprintf(w, `{"status":"ok","requests":"%v"}`, C.Get())
	}
}

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
