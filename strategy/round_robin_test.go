package strategy

import (
	"go-loadbalancer/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoundRobin_Distribution(t *testing.T) {
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
	for _, url := range urls {
		srv, err := server.NewBackendServer(url)
		if err != nil {
			t.Fatalf("could not create server: %v", err)
		}
		pool.AddServer(srv)
	}

	rr := NewRoundRobin()
	client := http.Client{}

	for i := 0; i < 99; i++ {
		target := rr.Next(pool)
		if target == nil {
			t.Fatal("no server returned")
		}
		resp, err := client.Get(target.GetUrl().String())
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
	}

	t.Logf("Request distribution:\nBackend 1: %d\nBackend 2: %d\nBackend 3: %d", counts[0], counts[1], counts[2])

	for i := range 3 {
		if abs(counts[i]-33) > 2 { // tolerate diff 2
			t.Errorf("Backend %d received %d requests; expected around 33", i+1, counts[i])
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func TestRoundRobin_AllServersDown(t *testing.T) {
	pool := server.NewPool()
	srv, _ := server.NewBackendServer("http://example.com")
	srv.SetAlive(false)
	pool.AddServer(srv)

	rr := NewRoundRobin()
	if rr.Next(pool) != nil {
		t.Error("Expected nil when all servers are down")
	}
}

func TestRoundRobin_SingleServerAlwaysSelected(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	pool := server.NewPool()
	s, _ := server.NewBackendServer(srv.URL)
	pool.AddServer(s)

	rr := NewRoundRobin()

	for range 5 {
		target := rr.Next(pool)
		if target == nil || target.GetUrl().String() != srv.URL {
			t.Errorf("expected the single server to always be returned, got: %v", target)
		}
	}
}
