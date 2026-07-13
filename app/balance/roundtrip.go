package balance

import (
	"github.com/vimcoders/grpcx"
	"google.golang.org/grpc"
)

type RoundTripper struct {
	sd *grpc.ServiceDesc
	grpcx.ClientConnInterface
}
