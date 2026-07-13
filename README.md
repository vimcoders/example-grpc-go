<div align="center">

<pre style="background:transparent">
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║   ██████╗  ██████╗ ██████╗  ██████╗ ██╗  ██╗               ║
║  ██╔════╝ ██╔═══██╗██╔══██╗██╔════╝ ╚██╗██╔╝               ║
║  ██║  ███╗██║   ██║██████╔╝██║  ███╗ ╚███╔╝                ║
║  ██║   ██║██║   ██║██╔══██╗██║   ██║ ██╔██╗                ║
║  ╚██████╔╝╚██████╔╝██║  ██║╚██████╔╝██╔╝ ██╗               ║
║   ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝               ║
║                                                               ║
║              ━━  Cloud-Native BFF Gateway  ━━               ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
</pre>

<p>
  <img src="https://img.shields.io/badge/Protocol-grpcx-00f0ff?style=for-the-badge&logo=go&logoColor=white&labelColor=0d1117" />
  <img src="https://img.shields.io/badge/Orchestrator-K8s-00f0ff?style=for-the-badge&logo=kubernetes&logoColor=white&labelColor=0d1117" />
  <img src="https://img.shields.io/badge/Performance-2.7x-00f0ff?style=for-the-badge&logo=sonic&logoColor=white&labelColor=0d1117" />
</p>

<p><i>Zero overhead. Maximum throughput. K8s native.</i></p>

</div>

---

## ⚡ Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              DATA FLOW                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   [Client]        [BFF Layer]         [Service Mesh]        [Microservice]  │
│      │                 │                      │                   │        │
│   HTTP│TLS│TCP    :26888              ClusterIP            :50051         │
│      │                 │                      │                   │        │
│      └───────────────▶│                      │                   │        │
│              Protocol   │───▶  Route  ───▶   │───▶  LoadBalance  │        │
│              Adapter    │      Match         │      HealthCheck   │        │
│                         │                      │                   │        │
│                         │◀─── Response ◀─────│◀─── grpcx ◀──────┘        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Design Principles**

| Layer | Responsibility | Technology |
|-------|---------------|------------|
| `Ingress` | Protocol adaptation, edge termination | TCP · HTTP · TLS |
| `BFF` | Request aggregation, transformation, caching | `app/balance` |
| `Mesh` | Service discovery, load balancing, health checks | K8s Native |
| `Transport` | High-performance RPC communication | [grpcx](https://github.com/vimcoders/grpcx) |

---

## 📊 Performance

```
╔═══════════════════════════════════════════════════════════════╗
║  BENCHMARK RESULTS                                            ║
╠═══════════════════════════════════════════════════════════════╣
║                                                               ║
║  grpcx (this stack)                                           ║
║  ─────────────────────────────────────────────                ║
║  ops/sec  │████████████████████████████████│  42,422          ║
║  latency  │█████████████░░░░░░░░░░░░░░░░░│  28,080 ns/op    ║
║  memory   │██████░░░░░░░░░░░░░░░░░░░░░░░░│  2,072 B/op      ║
║  allocs   │████░░░░░░░░░░░░░░░░░░░░░░░░░░│  44 allocs/op    ║
║                                                               ║
║  Standard gRPC                                                ║
║  ─────────────────────────────────────────────                ║
║  ops/sec  │████████████░░░░░░░░░░░░░░░░░░│  13,398          ║
║  latency  │█████████████████████████████│  80,192 ns/op    ║
║  memory   │████████████████████████░░░░░│  9,396 B/op      ║
║  allocs   │██████████████░░░░░░░░░░░░░░░│  153 allocs/op   ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝

Latency:     2.7x faster    ▲ 185%
Memory:      4.5x efficient ▲ 353%
Allocations: 3.5x fewer    ▲ 247%
```

---

## 🗂️ Structure

```
.
├── ⚙️ app/
│   ├── 🔷 balance/                 # BFF Gateway Core
│   │   ├── 📡 server.go            # TCP Listener & Connection Manager
│   │   ├── 🔀 session.go            # Per-Connection Request Router
│   │   ├── 📦 channel.go            # Binary Protocol (4-byte length + protobuf)
│   │   ├── 🔄 roundtrip.go          # Service Endpoint Wrapper
│   │   └── 🎛️ codec.go              # Protobuf Encoder/Decoder
│   └── 🟢 kube/                    # Backend Service Template
│       └── 📝 handler.go            # KubeService Implementation
│
├── 🚀 cmd/
│   ├── 🔷 balance/                 # Gateway Bootstrap
│   │   └── main.go
│   └── 🟢 kube/                    # Service Bootstrap
│       └── main.go
│
├── 📦 generated/
│   ├── 📋 proto/                   # IDL Definitions
│   └── 🔧 kubeapi/                 # Generated Go Stubs
│
├── 📈 kube_test.go                # Performance Regression Suite
├── 🐳 docker-compose.yaml         # Local Dev Environment
└── 📄 go.mod
```

---

## 🔧 Core Capabilities

### Multi-Service Routing

```go
// Initialize gateway
server := balance.NewServer()

// Register backend services via K8s DNS
server.RegisterService(&kubeapi.KubeService_ServiceDesc, "kube:50051")

// Start listening
server.ListenAndServe(ctx, ":26888")
```

**Routing Priority**
```
Priority 1  ──▶  Local Handlers        (e.g., HelloEcho)
Priority 2  ──▶  Service Descriptor    (grpc.ServiceDesc match)
Priority 3  ──▶  K8s Service Forward   (ClusterIP resolution)
```

### Protocol Stack

```
Layer 4  ┌─────────────────────────────────────┐
         │  Application                          │
         │  kubeapi.Request / kubeapi.Response   │
         ├─────────────────────────────────────┤
Layer 3  │  Serialization                        │
         │  protobuf (encoding/codec.go)          │
         ├─────────────────────────────────────┤
Layer 2  │  Framing                              │
         │  4-byte BE length + payload           │
         │  (app/balance/channel.go)            │
         ├─────────────────────────────────────┤
Layer 1  │  Transport                            │
         │  TCP / K8s ClusterIP                  │
         └─────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Local

```bash
# Terminal 1 ──▶ Backend
$ go run cmd/kube/main.go
[grpcx] listening on :50051

# Terminal 2 ──▶ Gateway
$ go run cmd/balance/main.go
[balance] listening on :26888

# Terminal 3 ──▶ Benchmark
$ go test -bench=. -benchmem -cpu=4,8,16
```

### Docker

```bash
$ docker-compose up -d
[+] Running 2/2
 ⠿ Container kube     Started
 ⠿ Container balance  Started
```

---

## 🗺️ Roadmap

```
[░░░░░░░░░░░░░░░░░░░░]  HTTP Ingress Adapter
[░░░░░░░░░░░░░░░░░░░░]  TLS Termination
[░░░░░░░░░░░░░░░░░░░░]  WebSocket Support
[░░░░░░░░░░░░░░░░░░░░]  Middleware Chain (Auth / RateLimit / Cache)
[░░░░░░░░░░░░░░░░░░░░]  Request Aggregation
[░░░░░░░░░░░░░░░░░░░░]  Distributed Tracing
```

---

## 🔗 Ecosystem

| Component | Repository | Role |
|-----------|------------|------|
| **grpcx** | [vimcoders/grpcx](https://github.com/vimcoders/grpcx) | High-performance RPC engine |
| **go-grpc-example** | [vimcoders/go-grpc-example](https://github.com/vimcoders/go-grpc-example) | BFF Gateway reference |

---

<div align="center">

**License** · Apache-2.0

</div>
