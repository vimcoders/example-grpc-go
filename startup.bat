@echo off
chcp 65001 >nul
cls

go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy,direct && go env -w GOSUMDB=sum.golang.google.cn
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
golangci-lint run --fix ./...

for /f "delims=" %%a in ('powershell Get-Date -Format "yyyy-MM-dd"') do set BUILD_DATE=%%a
for /f "delims=" %%i in ('git rev-parse --short HEAD') do set GIT_COMMIT=%%i
set IMAGE_ID=%BUILD_DATE%-%GIT_COMMIT%
docker compose up -d --build
docker exec -it redis-1 redis-cli --cluster create redis-1:6379 redis-2:6379 redis-3:6379 redis-4:6379 redis-5:6379 redis-6:6379 --cluster-replicas 1 --cluster-yes
docker system prune -af --filter "until=24h"
go test -bench . -cpu="1" -benchtime=1s -benchmem -count=1
echo ✅ 启动完成！
pause