# Golang HTTP Load Balancer

This is a modular, production-style HTTP Load Balancer built in **Go**, designed for learning and extensibility. It distributes incoming HTTP requests across multiple backend servers using customizable strategies.

## 🚀 Features

- 🔁 Round Robin and Weighted Round Robin strategies
- 🧱 Clean modular structure: server logic, pooling, and strategies are separated
- ✅ Thread-safe design using sync primitives
- 🔌 Real request forwarding using `httputil.ReverseProxy`
- 📦 Easily pluggable strategy architecture (Least Connections, IP Hash, etc.)
- 🧪 Unit-tested strategies with realistic scenarios using `httptest`

## 📁 Project Structure

```
go-loadbalancer/
├── main.go                      # Entry point of the application
├── server/
│   ├── server.go                # BackendServer and WeightedBackendServer definitions
│   └── pool.go                  # Server pool manager
├── strategy/
│   ├── round_robin.go           # Round Robin strategy
│   ├── weighted_round_robin.go  # Weighted Round Robin strategy
│   ├── round_robin_test.go      # Unit tests for Round Robin
│   └── weighted_round_robin_test.go # Unit tests for Weighted Round Robin
```

## 🧠 Components Overview

### `server/`
- **Server interface**: Defines the abstraction used by all load balancing strategies.
- **BackendServer**: A basic HTTP reverse proxy backend.
- **WeightedBackendServer**: Adds weight and current weight tracking for WRR.
- **Pool**: Thread-safe list of servers with add/get helpers.

### `strategy/`
- **RoundRobin**: Naive loop-based selection with `IsAlive()` check.
- **WeightedRoundRobin**: Smooth WRR implementation with dynamic weight balancing.
- Test cases cover single server, all-down, weight ratios, and edge conditions.

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
   http://localhost:8080/      # Round Robin
   http://localhost:8080/wrr   # Weighted Round Robin
   ```

## 🧪 Running Tests

To run all unit tests for Round Robin and Weighted Round Robin strategies:

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
```

## 📌 Notes

- All backends must respond with valid HTTP responses.
- Weighted strategies must only work with `WeightedBackendServer`.
- The `Next()` methods ensure servers are alive before selection.

## 📄 License

MIT License. Feel free to use, modify, and share!
