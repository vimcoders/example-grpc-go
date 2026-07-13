package main

import (
	"context"
	"example/app/balance"
	"example/generated/kubeapi"
)

func main() {
	server := balance.NewServer()
	if err := server.RegisterService(&kubeapi.KubeService_ServiceDesc, "kube:50051"); err != nil {
		panic(err)
	}
	server.ListenAndServe(context.Background(), ":26888")
}
