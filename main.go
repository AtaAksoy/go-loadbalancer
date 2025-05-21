package main

import (
	"fmt"
	"go-loadbalancer/server"
	"go-loadbalancer/strategy"
	"log"
	"net/http"
)

func main() {
	pool := server.NewPool()

	targets := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	for _, t := range targets {
		srv, err := server.NewBackendServer(t)
		if err != nil {
			log.Fatalf("Server could not added: %v", err)
		}

		pool.AddServer(srv)
	}

	rr := strategy.NewRoundRobin()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		target := rr.Next(pool)
		if target == nil {
			http.Error(w, "All servers down!", http.StatusServiceUnavailable)
			return
		}
		target.Serve(w, r)
	})

	fmt.Println("Load balancer runs on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
