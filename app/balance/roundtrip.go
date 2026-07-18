package balance

import (
	"context"

	"github.com/sony/gobreaker/v2"
	"github.com/vimcoders/grpcx"
	"github.com/vimcoders/grpcx/generated/api"
	"google.golang.org/grpc"
)

type RoundTripper struct {
	sd *grpc.ServiceDesc
	grpcx.ClientConnInterface
	cb *gobreaker.CircuitBreaker[*api.Response]
}

func (rt *RoundTripper) RoundTrip(ctx context.Context, req *api.Request) (*api.Response, error) {
	return rt.cb.Execute(func() (*api.Response, error) {
		return rt.ClientConnInterface.RoundTrip(ctx, req)
	})
}
