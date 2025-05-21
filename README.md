# Golang HTTP Load Balancer

This is a modular, production-style HTTP Load Balancer built in **Go**, designed for learning and extensibility. It distributes incoming HTTP requests across multiple backend servers using a **Round Robin** strategy by default.

## 🚀 Features

- 🔁 Round Robin request distribution strategy
- 🧱 Clean modular structure: server logic, pooling, and strategies are separated
- ✅ Thread-safe design using sync primitives
- 🔌 Real request forwarding using `httputil.ReverseProxy`
- 📦 Easily pluggable strategy architecture (Least Connections, IP Hash, etc.)

## 📁 Project Structure

```
loadbalancer/
├── main.go                    # Entry point of the application
├── server/
│   ├── server.go              # Backend server representation and proxy logic
│   └── pool.go                # Server pool manager
├── strategy/
│   └── round_robin.go         # Round Robin load balancing strategy
```

## 🧠 Components Overview

### `server/`
- **Server interface**: Defines the basic contract for backend servers.
- **BackendServer**: Implements reverse proxying, connection tracking, and health status.
- **Pool**: Manages the set of backend servers with thread safety.

### `strategy/`
- **RoundRobin**: Selects the next live server in a circular fashion.
- Future strategies like **Least Connections**, **IP Hashing** can be added easily.

## ▶️ Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/AtaAksoy/go-loadbalancer.git
   cd golang-loadbalancer
   ```

2. Add your backend servers to the `main.go` file:
   ```go
   targets := []string{
       "http://localhost:8081",
       "http://localhost:8082",
   }
   ```

3. Run the load balancer:
   ```bash
   go run main.go
   ```

4. Send requests to:
   ```
   http://localhost:8080
   ```

## 📌 Notes

- Backend servers must be running and return valid HTTP responses.
- This project is intended as a learning tool and can be expanded for production use with features like:
  - Health checks
  - Retry policies
  - Request logging
  - HTTPS support

## 📄 License

MIT License. Feel free to use, modify, and share!
