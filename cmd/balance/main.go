package main

import (
	"context"
	"example/app/balance"
	"example/generated/kubeapi"
	"math"
	"net/http"
	"time"
)

func main() {
	server := balance.NewServer()
	if err := server.RegisterService(&kubeapi.KubeService_ServiceDesc, "kube:50051"); err != nil {
		panic(err)
	}
	server.ListenAndServe(context.Background(), ":26888")
	svr := &http.Server{
		// 监听地址，格式：:端口号
		Addr: ":26888",
		// 路由处理器（使用当前 Handler 注册的路由）
		Handler: server,
		// 读取请求超时时间：10秒
		ReadTimeout: 10 * time.Second,
		// 响应写入超时时间：10秒
		WriteTimeout: 10 * time.Second,
		// 长连接空闲超时时间：30秒
		IdleTimeout: 30 * time.Second,
		// 最大请求头大小（使用 math.MaxInt16 避免溢出）
		MaxHeaderBytes: math.MaxInt16,
	}
	svr.ListenAndServe()
}
