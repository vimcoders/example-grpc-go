package main

import (
	"context"
	"example/app/activity"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.ActivityService_ServiceDesc, &activity.Handler{})
	server.ListenAndServe(context.Background(), ":50054")
}
