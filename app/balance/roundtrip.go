package balance

import (
	"context"

	"github.com/vimcoders/grpcx"
	"github.com/vimcoders/grpcx/generated/api"
	"google.golang.org/grpc"
)

type RoundTripper struct {
	grpc.ServiceDesc
	grpcx.ClientConnInterface
}

func (rt *RoundTripper) RoundTrip(ctx context.Context, req *api.Request) (*api.Response, error) {
	// TODO:: 熔断器
	return rt.ClientConnInterface.RoundTrip(ctx, req)
}

func (rt *RoundTripper) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	// TODO:: 熔断器
	return rt.ClientConnInterface.Invoke(ctx, method, args, reply, opts...)
}
