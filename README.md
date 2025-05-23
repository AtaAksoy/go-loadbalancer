# Golang HTTP Load Balancer

This is a modular, production-style HTTP Load Balancer built in **Go**, designed for learning and extensibility. It distributes incoming HTTP requests across multiple backend servers using customizable strategies.

## 🚀 Features

- 🔁 Round Robin, Weighted Round Robin, and Least Connections strategies
- 🧱 Clean modular structure: server logic, pooling, and strategies are separated
- ✅ Thread-safe design using `sync.Mutex` and `sync.RWMutex`
- 🔌 Real request forwarding using `httputil.ReverseProxy`
- 🧪 Realistic, concurrency-aware unit tests using `httptest` and `goroutines`
- 📦 Easily extensible strategy layer for future additions (e.g., IP Hashing, Health Checks)

## 📁 Project Structure

```
go-loadbalancer/
├── main.go                         # Entry point of the application
├── server/
│   ├── server.go                   # BackendServer and WeightedBackendServer definitions
│   └── pool.go                     # Server pool manager
├── strategy/
│   ├── round_robin.go              # Round Robin strategy
│   ├── weighted_round_robin.go     # Weighted Round Robin strategy
│   ├── least_connection.go         # Least Connection strategy
│   ├── round_robin_test.go         # Unit tests for Round Robin
│   ├── weighted_round_robin_test.go # Unit tests for Weighted Round Robin
│   └── least_connection_test.go    # Unit tests for Least Connection
```

## 🧠 Components Overview

### `server/`
- **Server interface**: Defines the abstraction used by all strategies.
- **BackendServer**: Basic reverse proxy logic with connection counting.
- **WeightedBackendServer**: Adds weight/current weight tracking for WRR.
- **Pool**: Stores a set of backend servers in a thread-safe list.

### `strategy/`
- **RoundRobin**: Rotates through live servers in sequence.
- **WeightedRoundRobin**: Prioritizes servers based on configured weights, using smooth WRR logic.
- **LeastConnection**: Selects the server with the fewest active connections in real-time.

## ▶️ Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/AtaAksoy/go-loadbalancer.git
   cd go-loadbalancer
   ```

2. Add your backend servers to the `main.go` file:
   ```go
   targets := []string{
       "http://localhost:8081",
       "http://localhost:8082",
       "http://localhost:8083",
   }
   ```

3. Run the load balancer:
   ```bash
   go run main.go
   ```

4. Send requests to:
   ```
   http://localhost:8080/       # Round Robin
   http://localhost:8080/wrr    # Weighted Round Robin
   http://localhost:8080/least  # Least Connection
   ```

## 🧪 Running Tests

Run all unit tests for all strategies:
```bash
go test ./strategy -v
```

### Sample Test Output

```
=== RUN   TestWeightedRoundRobin_Distribution
    weighted_round_robin_test.go:51: Request distribution:
        Backend 1: 50
        Backend 2: 30
        Backend 3: 20

=== RUN   TestLeastConnection_WithConcurrentServe
    least_connection_test.go:72: Request distribution with concurrent Serve(): [32 31 27]
```

## 📌 Notes

- Make sure backends are running and accessible before testing manually.
- `Serve()` must be used in tests to correctly simulate active connections.
- All strategies skip servers marked as down via `IsAlive()`.

## 📄 License

MIT License. Feel free to use, modify, and share!
