package server

import "sync"

type ServerPool interface {
	AddServer(Server)
	GetServers() []Server
	GetServerPoolSize() int
}

type Pool struct {
	servers []Server
	mu      sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		servers: []Server{},
	}
}

func (p *Pool) AddServer(s Server) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers = append(p.servers, s)
}

func (p *Pool) GetServers() []Server {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.servers
}

func (p *Pool) GetServerPoolSize() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.servers)
}
