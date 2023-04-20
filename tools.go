package main

import (
	"log"
	"os"
	"strconv"
	"sync"
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
