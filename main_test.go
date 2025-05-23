package main

import (
	"go-loadbalancer/server"
	"go-loadbalancer/strategy"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeightedRoundRobin(t *testing.T) {
	counts := make([]int, 3)

	handlers := []http.HandlerFunc{
		func(w http.ResponseWriter, r *http.Request) { counts[0]++ },
		func(w http.ResponseWriter, r *http.Request) { counts[1]++ },
		func(w http.ResponseWriter, r *http.Request) { counts[2]++ },
	}

	var urls []string
	for _, h := range handlers {
		ts := httptest.NewServer(http.HandlerFunc(h))
		defer ts.Close()
		urls = append(urls, ts.URL)
	}

	pool := server.NewPool()
	weights := []int{5, 3, 2}
	for i, url := range urls {
		srv, _ := server.NewWeightedBackendServer(url, weights[i])
		pool.AddServer(srv)
	}

	wrr := strategy.NewWeightedRoundRobin()
	client := http.Client{}

	for range 105 {
		target := wrr.Next(pool)
		resp, err := client.Get(target.GetUrl().String())
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
	}

	t.Logf("Backend distribution: %v", counts)
}
