package strategy

import (
	"go-loadbalancer/server"
	"sync"
)

type RoundRobin struct {
	counter int
	mu      sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (rr *RoundRobin) Next(p *server.Pool) server.Server {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	servers := p.GetServers()
	n := len(servers)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		s := servers[rr.counter%n]
		rr.counter = (rr.counter + 1) % n
		if s.IsAlive() {
			return s
		}
	}

	return nil
}
