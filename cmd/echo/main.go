package main

import (
	"context"
	"kube/app/echo"
	"kube/generated/api"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&api.EchoService_ServiceDesc, &echo.Handler{})
	server.ListenAndServe(context.Background(), ":50051")
}
