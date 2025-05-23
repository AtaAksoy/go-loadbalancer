package main

import (
	"fmt"
	"go-loadbalancer/server"
	"go-loadbalancer/strategy"
	"log"
	"net/http"
)

func main() {
	testWeightedRoundRobin()
}

func testRoundRobin() {
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

func testWeightedRoundRobin() {
	pool := server.NewPool()

	servers := []struct {
		URL    string
		Weight int
	}{
		{"http://localhost:8081", 5},
		{"http://localhost:8082", 2},
		{"http://localhost:8083", 1},
	}

	for _, s := range servers {
		srv, err := server.NewWeightedBackendServer(s.URL, s.Weight)
		if err != nil {
			log.Fatalf("Weighted server could not be added: %v", err)
		}
		pool.AddServer(srv)
	}

	wrr := strategy.NewWeightedRoundRobin()

	http.HandleFunc("/wrr", func(w http.ResponseWriter, r *http.Request) {
		target := wrr.Next(pool)
		if target == nil {
			http.Error(w, "All servers down!", http.StatusServiceUnavailable)
			return
		}
		target.Serve(w, r)
	})
}
