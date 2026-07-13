# go-grpc-example

<p align="center">
  <img src="https://img.shields.io/badge/grpcx-powered-blue?style=flat-square&logo=go" />
  <img src="https://img.shields.io/badge/k8s-native-green?style=flat-square&logo=kubernetes" />
  <img src="https://img.shields.io/badge/latency-2.7x%20faster-orange?style=flat-square" />
</p>

> 面向 K8s 云原生环境的高性能 BFF 网关，基于自研 grpcx 协议栈构建。

## 架构概览

```
┌──────────────┐      ┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│   Client     │─────▶│  BFF Layer   │─────▶│  K8s Mesh   │─────▶│  Microservice│
│  HTTP/TCP    │      │  :26888      │      │  ClusterIP   │      │  :50051      │
│  WebSocket   │      │              │      │  Service     │      │  grpcx       │
└──────────────┘      └──────────────┘      └──────────────┘      └──────────────┘
       │                     │                     │                     │
       ▼                     ▼                     ▼                     ▼
  Protocol              Service                Load                  Business
  Adapter               Router                 Balance               Logic
```

### 设计哲学

| 层级 | 职责 | 技术选型 |
|------|------|----------|
| **Ingress** | 协议适配、边缘接入 | TCP / HTTP / TLS |
| **BFF** | 请求聚合、协议转换、缓存 | `app/balance` |
| **Service Mesh** | 服务发现、负载均衡、健康检查 | K8s Native |
| **Transport** | 高性能 RPC 通信 | [grpcx](https://github.com/vimcoders/grpcx) |

## 性能基准

```
BenchmarkEcho-4        42422    28080 ns/op    2072 B/op    44 allocs/op
BenchmarkStdGRPC-4     13398    80192 ns/op    9396 B/op   153 allocs/op

Latency:    2.7x faster
Memory:     4.5x efficient
Allocations: 3.5x fewer
```

## 项目结构

```
.
├── app/
│   ├── balance/          # BFF Gateway Core
│   │   ├── server.go     # TCP Listener & Connection Manager
│   │   ├── session.go     # Per-Connection Request Router
│   │   ├── channel.go     # Binary Protocol (4-byte length + protobuf)
│   │   ├── roundtrip.go   # Service Endpoint Wrapper
│   │   └── codec.go       # Protobuf Encoder/Decoder
│   └── kube/              # Backend Service Template
│       └── handler.go     # KubeService Implementation
├── cmd/
│   ├── balance/           # Gateway Bootstrap
│   │   └── main.go
│   └── kube/              # Service Bootstrap
│       └── main.go
├── generated/
│   ├── proto/             # IDL Definitions
│   └── kubeapi/           # Generated Go Stubs
├── kube_test.go           # Performance Regression Suite
├── docker-compose.yaml    # Local Dev Environment
└── go.mod
```

## 核心能力

### Multi-Service Routing

```go
server := balance.NewServer()
server.RegisterService(&kubeapi.KubeService_ServiceDesc, "kube:50051")
server.ListenAndServe(ctx, ":26888")
```

**Routing Priority:**
1. Local Handlers (e.g., `HelloEcho`)
2. Service Descriptor Match (`grpc.ServiceDesc`)
3. K8s Service Forwarding

### Protocol Stack

```
┌─────────────────────────────────────┐
│  Application Layer                  │
│  kubeapi.Request / kubeapi.Response   │
├─────────────────────────────────────┤
│  Serialization Layer                │
│  protobuf (encoding/codec.go)        │
├─────────────────────────────────────┤
│  Framing Layer                      │
│  4-byte BE length + payload          │
│  (app/balance/channel.go)           │
├─────────────────────────────────────┤
│  Transport Layer                      │
│  TCP / K8s ClusterIP                │
└─────────────────────────────────────┘
```

## Quick Start

### Local Development

```bash
# Terminal 1: Start Backend
$ go run cmd/kube/main.go
[grpcx] listening on :50051

# Terminal 2: Start Gateway
$ go run cmd/balance/main.go
[balance] listening on :26888

# Terminal 3: Benchmark
$ go test -bench=. -benchmem -cpu=4,8,16
```

### Docker Compose

```bash
$ docker-compose up -d
[+] Running 2/2
 ⠿ Container kube     Started
 ⠿ Container balance  Started
```

## Roadmap

- [ ] HTTP Ingress Adapter
- [ ] TLS Termination
- [ ] WebSocket Support
- [ ] Middleware Chain (Auth / RateLimit / Cache)
- [ ] Request Aggregation
- [ ] Distributed Tracing

## Ecosystem

| Component | Repository | Role |
|-----------|------------|------|
| grpcx | [vimcoders/grpcx](https://github.com/vimcoders/grpcx) | High-performance RPC engine |
| go-grpc-example | [vimcoders/go-grpc-example](https://github.com/vimcoders/go-grpc-example) | BFF Gateway reference implementation |

## License

Apache-2.0
