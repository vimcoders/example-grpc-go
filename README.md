# go-grpc-example

基于 [grpcx](https://github.com/vimcoders/grpcx) 的 BFF（Backend for Frontend）示例项目，展示如何在 K8s 环境中构建高性能的微服务网关。

## 架构设计

```
┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Client │────>│   BFF       │────>│  K8s SVC    │────>│  Backend    │
│ (HTTP)  │     │  (balance)  │     │  (ClusterIP)│     │  (grpcx)    │
│ (TLS)   │     │  (:26888)   │     │  user-svc   │     │  (:50051)   │
│ (TCP)   │     │             │     │  order-svc  │     │             │
└─────────┘     └─────────────┘     └─────────────┘     └─────────────┘
```

**设计原则：**
- **K8s 负责基础设施**：服务发现、负载均衡、健康检查
- **BFF 负责业务适配**：协议转换、请求聚合、鉴权缓存
- **grpcx 负责服务通信**：高性能 RPC，比标准 gRPC 快 2.7x

## 项目结构

```
.
├── app/
│   ├── balance/          # BFF 网关层
│   │   ├── server.go     # 网关入口，监听 TCP
│   │   ├── session.go    # 连接会话，处理请求路由
│   │   ├── channel.go    # 自定义协议（4 字节 length + payload）
│   │   ├── roundtrip.go  # RoundTripper 包装
│   │   └── codec.go      # protobuf 编解码
│   └── kube/             # 后端服务示例
│       └── handler.go    # KubeService 实现
├── cmd/
│   ├── balance/          # BFF 启动入口
│   │   └── main.go
│   └── kube/             # 后端服务启动入口
│       └── main.go
├── generated/            # protoc 生成代码
│   ├── proto/            # .proto 定义
│   └── kubeapi/          # Go 生成代码
├── kube_test.go          # 基准测试
├── docker-compose.yaml   # 本地开发环境
└── go.mod
```

## 核心组件

### 1. BFF 网关 (app/balance)

**多服务路由：**

```go
server := balance.NewServer()
server.RegisterService(&kubeapi.KubeService_ServiceDesc, "kube:50051")
server.ListenAndServe(context.Background(), ":26888")
```

**路由逻辑：**
1. 本地方法优先（如 HelloEcho）
2. 匹配 endpoints 中的 ServiceDesc
3. 转发到对应 K8s Service

**协议格式：**
- 4 字节大端 length
- protobuf payload

### 2. 后端服务 (app/kube)

```go
server := grpcx.NewServer()
server.RegisterService(&kubeapi.KubeService_ServiceDesc, &kube.Handler{})
server.ListenAndServe(context.Background(), ":50051")
```

## 快速开始

### 本地开发

```bash
# 启动后端服务
go run cmd/kube/main.go

# 启动 BFF 网关
go run cmd/balance/main.go

# 运行基准测试
go test -bench=BenchmarkHello -benchmem
```

### Docker Compose

```bash
docker-compose up -d
```

## 基准测试

```bash
go test -bench=. -benchmem -cpu=4,8,16
```

| 指标 | 说明 |
|------|------|
| ns/op | 每次操作耗时 |
| B/op | 每次操作分配字节数 |
| allocs/op | 每次操作堆分配次数 |

## 扩展计划

- [ ] HTTP 接入层
- [ ] TLS 终止
- [ ] WebSocket 支持
- [ ] 中间件机制（鉴权、日志、限流）
- [ ] 请求聚合（减少前端请求数）

## 依赖

- [grpcx](https://github.com/vimcoders/grpcx) - 轻量级 RPC 框架
- [protobuf](https://github.com/protocolbuffers/protobuf) - 序列化
- [grpc-go](https://github.com/grpc/grpc-go) - 兼容 gRPC API

## License

Apache-2.0
