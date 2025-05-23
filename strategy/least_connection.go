package strategy

import (
	"go-loadbalancer/server"
	"sync"
)

type LeastConnection struct {
	mu sync.Mutex
}

func NewLeastConnection() *LeastConnection {
	return &LeastConnection{}
}

func (lc *LeastConnection) Next(p *server.Pool) server.Server {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	servers := p.GetServers()
	if n := p.GetServerPoolSize(); n == 0 {
		return nil
	}
	minConnections := -1
	var selected server.Server

	for _, s := range servers {
		if !s.IsAlive() {
			continue
		}

		conn := s.GetActiveConnections()
		if selected == nil || conn < minConnections {
			selected = s
			minConnections = conn
		}
	}

	return selected
}
