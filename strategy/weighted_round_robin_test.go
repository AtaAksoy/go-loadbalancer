package strategy

import (
	"go-loadbalancer/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeightedRoundRobin_Distribution(t *testing.T) {
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
		srv, err := server.NewWeightedBackendServer(url, weights[i])
		if err != nil {
			t.Fatalf("failed to create server: %v", err)
		}
		pool.AddServer(srv)
	}

	wrr := NewWeightedRoundRobin()
	client := http.Client{}

	for i := 0; i < 100; i++ {
		target := wrr.Next(pool)
		if target == nil {
			t.Fatal("no server returned by strategy")
		}
		resp, err := client.Get(target.GetUrl().String())
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
	}

	t.Logf("Request distribution:\nBackend 1: %d\nBackend 2: %d\nBackend 3: %d", counts[0], counts[1], counts[2])

	expected := []int{50, 30, 20}
	tolerance := 10

	for i, count := range counts {
		diff := expected[i] - count
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			t.Errorf("Backend %d out of expected range: got %d, expected ~%d", i+1, count, expected[i])
		}
	}
}

func TestWeightedRoundRobin_SingleServer(t *testing.T) {
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
	}))
	defer ts.Close()

	pool := server.NewPool()
	srv, err := server.NewWeightedBackendServer(ts.URL, 10)
	if err != nil {
		t.Fatalf("could not create server: %v", err)
	}
	pool.AddServer(srv)

	wrr := NewWeightedRoundRobin()
	client := http.Client{}

	for i := 0; i < 50; i++ {
		target := wrr.Next(pool)
		if target == nil {
			t.Fatal("expected a server but got nil")
		}
		resp, err := client.Get(target.GetUrl().String())
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
	}

	if count != 50 {
		t.Errorf("expected all 50 requests to go to the single server, got %d", count)
	}
}

func TestWeightedRoundRobin_AllServersDown(t *testing.T) {
	pool := server.NewPool()
	srv, _ := server.NewWeightedBackendServer("http://fake", 5)
	srv.SetAlive(false)
	pool.AddServer(srv)

	wrr := NewWeightedRoundRobin()
	target := wrr.Next(pool)
	if target != nil {
		t.Error("expected nil when all servers are down")
	}
}

func TestWeightedRoundRobin_ZeroOrLowWeights(t *testing.T) {
	counts := make([]int, 2)

	handlers := []http.HandlerFunc{
		func(w http.ResponseWriter, r *http.Request) { counts[0]++ },
		func(w http.ResponseWriter, r *http.Request) { counts[1]++ },
	}

	var urls []string
	for _, h := range handlers {
		ts := httptest.NewServer(http.HandlerFunc(h))
		defer ts.Close()
		urls = append(urls, ts.URL)
	}

	pool := server.NewPool()
	weights := []int{1, 0}
	for i, url := range urls {
		srv, _ := server.NewWeightedBackendServer(url, weights[i])
		pool.AddServer(srv)
	}

	wrr := NewWeightedRoundRobin()
	client := http.Client{}

	for i := 0; i < 20; i++ {
		target := wrr.Next(pool)
		if target == nil {
			t.Fatal("no server returned")
		}
		resp, err := client.Get(target.GetUrl().String())
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
	}

	if counts[0] < 18 || counts[1] > 2 {
		t.Errorf("unexpected distribution: %+v", counts)
	}
}
