package main

import (
	"context"
	"example/app/social"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.SocialService_ServiceDesc, &social.Handler{})
	server.ListenAndServe(context.Background(), ":50052")
}
