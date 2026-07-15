package main

import (
	"context"
	"example/app/mail"
	"example/generated/kubeapi"

	"github.com/vimcoders/grpcx"
)

func main() {
	server := grpcx.NewServer()
	server.RegisterService(&kubeapi.MailService_ServiceDesc, &mail.Handler{})
	server.ListenAndServe(context.Background(), ":50056")
}
