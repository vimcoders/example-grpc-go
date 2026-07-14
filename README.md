# example-grpc-go

使用 [grpcx](https://github.com/vimcoders/grpcx) 构建 **微服务网关路由** 的示例项目。

### 环境要求

- Go 1.26+
- Docker & Docker Compose

### 1. 克隆仓库

```bash
git clone https://github.com/vimcoders/example-grpc-go.git
cd example-grpc-go
```

### 2. 启动服务

```bash
docker compose up -d
```

启动的服务：
- `balance` 网关，端口 `:26888`
- `kube` 微服务，端口 `:50051`
- `mysql`、`redis`、`nats` 基础设施

### 3. 压测

```bash
go test -bench=. -v ./...
```

## 项目结构

```
.
├── app/
│   ├── balance/          # BFF 网关实现
│   │   ├── server.go     # TCP 监听器
│   │   ├── session.go    # 请求路由分发
│   │   ├── roundtrip.go  # 后端端点封装
│   │   ├── channel.go    # 帧协议编解码
│   │   └── codec.go      # Protobuf 编解码器
│   └── kube/
│       └── handler.go    # 业务逻辑处理器
├── cmd/
│   ├── balance/          # 网关入口
│   │   ├── main.go
│   │   └── Dockerfile
│   └── kube/             # 微服务入口
│       ├── main.go
│       └── Dockerfile
├── generated/
│   ├── proto/            # .proto 定义文件
│   └── kubeapi/          # 生成的 Go 代码
├── docker-compose.yaml   # 本地开发环境
├── kube_test.go          # 压测客户端
└── go.mod
```

## 工作原理

### balance（网关）

注册路由监听 `:26888`，通过 grpcx 转发到后端服务：

```go
server := balance.NewServer()
server.RegisterService(&kubeapi.KubeService_ServiceDesc, "kube:50051")
server.ListenAndServe(ctx, ":26888")
```

### kube（微服务）

标准 grpcx 服务：

```go
server := grpcx.NewServer()
server.RegisterService(&kubeapi.KubeService_ServiceDesc, &kube.Handler{})
server.ListenAndServe(ctx, ":50051")
```

参考客户端实现见 `kube_test.go`。

## 依赖

- [grpcx](https://github.com/vimcoders/grpcx)
- [google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc)
- [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf)

## 许可证

Apache-2.0
