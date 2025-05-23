package strategy

import (
	"go-loadbalancer/server"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestLeastConnection_WithConcurrentServe(t *testing.T) {
	pool := server.NewPool()
	counts := make([]int, 1000)
	var servers []server.Server

	for i := range len(counts) {
		idx := i
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(50 * time.Millisecond)
			counts[idx]++
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		srv, err := server.NewBackendServer(ts.URL)
		if err != nil {
			t.Fatalf("failed to create server: %v", err)
		}
		pool.AddServer(srv)
		servers = append(servers, srv)
	}

	lc := NewLeastConnection()
	var wg sync.WaitGroup
	req, _ := http.NewRequest("GET", "/", nil)

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			target := lc.Next(pool)
			if target == nil {
				t.Error("no server returned")
				return
			}
			w := httptest.NewRecorder()
			target.Serve(w, req)
		}()
	}

	wg.Wait()
	t.Logf("Request distribution with concurrent Serve(): %v", counts)
}

func TestLeastConnection_SingleServer(t *testing.T) {
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
	}))
	defer ts.Close()

	pool := server.NewPool()
	srv, _ := server.NewBackendServer(ts.URL)
	pool.AddServer(srv)

	lc := NewLeastConnection()
	client := http.Client{}

	for i := 0; i < 20; i++ {
		target := lc.Next(pool)
		if target == nil {
			t.Fatal("expected a server")
		}
		resp, err := client.Get(target.GetUrl().String())
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
	}

	if count != 20 {
		t.Errorf("expected 20 requests to be handled, got %d", count)
	}
}

func TestLeastConnection_AllServersDown(t *testing.T) {
	pool := server.NewPool()
	srv, _ := server.NewBackendServer("http://example.com")
	srv.SetAlive(false)
	pool.AddServer(srv)

	lc := NewLeastConnection()
	target := lc.Next(pool)
	if target != nil {
		t.Error("expected nil when all servers are down")
	}
}
