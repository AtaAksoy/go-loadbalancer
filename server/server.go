package server

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server interface {
	SetAlive(bool)
	IsAlive() bool
	GetUrl() *url.URL
	GetActiveConnections() int
	Serve(http.ResponseWriter, *http.Request)
}

type BackendServer struct {
	URL          *url.URL
	alive        bool
	mu           sync.RWMutex
	connections  int
	reverseProxy *httputil.ReverseProxy
}

type WeightedBackendServer struct {
	*BackendServer
	Weight        int
	CurrentWeight int
}

func NewBackendServer(target string) (*BackendServer, error) {
	urlParsed, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &BackendServer{
		URL:          urlParsed,
		alive:        true,
		reverseProxy: httputil.NewSingleHostReverseProxy(urlParsed),
	}, nil
}

func NewWeightedBackendServer(target string, weight int) (*WeightedBackendServer, error) {
	base, err := NewBackendServer(target)
	if err != nil {
		return nil, err
	}
	return &WeightedBackendServer{
		BackendServer: base,
		Weight:        weight,
		CurrentWeight: 0,
	}, nil
}

func (s *BackendServer) SetAlive(alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alive = alive
}

func (s *BackendServer) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.alive
}

func (s *BackendServer) GetUrl() *url.URL {
	return s.URL
}

func (s *BackendServer) GetActiveConnections() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connections
}

func (s *BackendServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.connections++
	s.mu.Unlock()

	log.Printf("🔥 Request redirects to: %v", s.GetUrl())
	s.reverseProxy.ServeHTTP(w, r)

	s.mu.Lock()
	s.connections--
	s.mu.Unlock()
}
