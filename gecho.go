package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
