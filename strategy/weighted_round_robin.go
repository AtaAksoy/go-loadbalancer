package strategy

import (
	"go-loadbalancer/server"
	"sync"
)

type WeightedRoundRobin struct {
	mu sync.Mutex
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{}
}

func (wrr *WeightedRoundRobin) Next(p *server.Pool) server.Server {
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	servers := p.GetServers()
	n := p.GetServerPoolSize()
	if n == 0 {
		return nil
	}

	var (
		totalWeight = 0
		maxWeight   = 0
		best        server.Server
	)

	for _, s := range servers {
		if !s.IsAlive() {
			return nil
		}
		ws, ok := s.(*server.WeightedBackendServer)
		if !ok {
			continue
		}

		ws.CurrentWeight += ws.Weight
		totalWeight += ws.Weight

		if ws.CurrentWeight > maxWeight {
			maxWeight = ws.CurrentWeight
			best = ws
		}
	}

	if best != nil {
		best.(*server.WeightedBackendServer).CurrentWeight -= totalWeight
	}

	return best
}
