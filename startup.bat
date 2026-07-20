@echo off
chcp 65001 >nul
cls

echo 🎨 统一代码换行符配置
git config core.autocrlf input
echo 🚀 配置Golang国内代理源
go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy,direct && go env -w GOSUMDB=sum.golang.google.cn
echo 💡 安装代码自动补全 gopls
go install golang.org/x/tools/gopls@latest
echo 🧹 安装代码规范检测 golangci-lint v2
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
echo 📦 安装 gRPC 代码生成工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo ✨ 安装Protobuf格式化工具
go install github.com/bufbuild/buf/cmd/buf@latest
echo 🎨 执行代码自动格式化与规范检测
golangci-lint run --fix ./...

for /f "delims=" %%a in ('powershell Get-Date -Format "yyyy-MM-dd"') do set BUILD_DATE=%%a
for /f "delims=" %%i in ('git rev-parse --short HEAD') do set GIT_COMMIT=%%i
set IMAGE_ID=%BUILD_DATE%-%GIT_COMMIT%
echo 🐳 启动中间件集群(Mysql/Redis/Nats)
docker compose up -d mysql redis-1 redis-2 redis-3 redis-4 redis-5 redis-6 nats-1 nats-2 nats-3 redis-cluster-init
echo 🐳 构建并启动业务服务容器
docker compose up -d --build
docker system prune -af --filter "until=24h"
echo ⚡ 执行基准性能压测
go test ./test/bench -bench . -cpu="1" -benchtime=1s -benchmem -count=1
echo ✅ 全部开发环境初始化完毕
pause