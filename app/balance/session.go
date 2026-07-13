package balance

import (
	"context"
	"kube/generated/api"
	"net"
	"path"
	"slices"
	"time"

	"github.com/vimcoders/grpcx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

type Session struct {
	encoding.Codec
	desc        *grpc.ServiceDesc
	interceptor grpc.UnaryServerInterceptor
	endpoints   []grpcx.ClientConnInterface
}

func (s *Session) Echo(ctx context.Context, req *api.EchoRequest) (*api.EchoResponse, error) {
	return &api.EchoResponse{Message: req.Message}, nil
}

func (s *Session) RoundTrip(ctx context.Context, req *api.Request) (*api.Response, error) {
	if req.Method == "" {
		return &api.Response{
			Code:    int32(codes.OK),
			Message: codes.OK.String(),
		}, nil
	}
	idx := slices.IndexFunc(s.desc.Methods, func(v grpc.MethodDesc) bool {
		return path.Join("/", s.desc.ServiceName, v.MethodName) == req.Method
	})
	if idx < 0 {
		return &api.Response{
			Code:    int32(codes.Unimplemented),
			Message: codes.Unimplemented.String(),
		}, nil
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(req.Timeout)*time.Millisecond)
	defer cancel()
	reply, err := s.desc.Methods[idx].Handler(
		s,
		timeoutCtx,
		func(in any) error {
			return s.Unmarshal(req.Payload, in)
		},
		s.interceptor)
	if err != nil {
		return &api.Response{
			Code:    int32(codes.Unavailable),
			Message: err.Error(),
		}, nil
	}
	response, err := s.Marshal(reply)
	if err != nil {
		return &api.Response{
			Code:    int32(codes.Unavailable),
			Message: err.Error(),
		}, nil
	}
	return &api.Response{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Payload: response,
	}, nil
}

func (s *Session) Handle(ctx context.Context, c net.Conn) (err error) {
	defer c.Close()
	channel := newChannel(c)
	for {
		select {
		case <-ctx.Done():
			return status.Error(codes.Canceled, codes.Canceled.String())
		default:
			payload, err := channel.Recv()
			if err != nil {
				return err
			}
			var request api.Request
			if err := s.Unmarshal(payload, &request); err != nil {
				return err
			}
			channel.putmbuf(payload)
			response, err := s.RoundTrip(ctx, &request)
			if err != nil {
				return err
			}
			b, err := s.Marshal(response)
			if err != nil {
				return err
			}
			if err := channel.Send(b); err != nil {
				return err
			}
		}
	}
}
