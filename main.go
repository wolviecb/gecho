package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Counter mutex for counter
type counter struct {
	mu sync.Mutex
	n  int
}

// Global counter
var C counter

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
			log.Fatalf("Invalid Value, %s not a valid integer", k)
		}
		r = i
	default:
		log.Fatalf("Invalid Value, %T not valid", k)
	}
	return r.(D)
}

// serve evaluates if the limit of requests is reached
func serve(l int) bool {
	return C.Get() < l
}

// reverse returns a reversed byte array of c
func reverse[C ~[]E, E any](c C) C {
	for i, j := 0, len(c)-1; i < j; i, j = i+1, j-1 {
		c[i], c[j] = c[j], c[i]
	}
	return c
}

// handler generates the echo server response
func handler(l int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer C.Add()
		if !serve(l) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Header.Get("Content-Type") == "application/json" {
			v := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				http.Error(w, "why did you do that?", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "%s", v)
			return
		}
		v, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
		}
		if r.Header.Get("Reverse") == "true" {
			v = reverse(v)
		}
		fmt.Fprintf(w, "%s", v)
	}
}

// httpHealth is a bad implementation of a health check
func httpHealth(l int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !serve(l) {
			http.Error(w, fmt.Sprintf(`{"status":"FAIL","requests":"%v"}`, C.Get()), http.StatusInternalServerError)
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
			http.Error(w, "Bad request, invalid token", http.StatusBadRequest)
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
