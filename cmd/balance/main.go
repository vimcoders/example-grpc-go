package main

import (
	"context"
	"kube/app/balance"
	"kube/generated/kubeapi"
)

func main() {
	server := balance.NewServer()
	server.RegisterService(&kubeapi.HelloService_ServiceDesc, ":50051")
	server.ListenAndServe(context.Background(), ":26888")
}
