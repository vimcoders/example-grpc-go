package balance

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
	"net"
	"path"
	"slices"
	"time"

	"github.com/vimcoders/grpcx/generated/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"
)

type Session struct {
	kubeapi.HelloServiceServer
	encoding.Codec
	desc        *grpc.ServiceDesc
	interceptor grpc.UnaryServerInterceptor
	endpoints   []RoundTripper
}

func (s *Session) Hello(ctx context.Context, req *kubeapi.HelloRequest) (*kubeapi.HelloResponse, error) {
	slog.Info("echo", "hello", req)
	return &kubeapi.HelloResponse{Message: req.Message}, nil
}

func (s *Session) RoundTrip(ctx context.Context, req *kubeapi.Request) (*kubeapi.Response, error) {
	if req.Method == "" {
		return &kubeapi.Response{
			Code:    int32(codes.OK),
			Message: codes.OK.String(),
		}, nil
	}
	idx := slices.IndexFunc(s.desc.Methods, func(v grpc.MethodDesc) bool {
		return v.MethodName == req.Method
	})
	if idx >= 0 {
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
			return &kubeapi.Response{
				Code:    int32(codes.Unavailable),
				Message: err.Error(),
			}, nil
		}
		response, err := s.Marshal(reply)
		if err != nil {
			return &kubeapi.Response{
				Code:    int32(codes.Unavailable),
				Message: err.Error(),
			}, nil
		}
		return &kubeapi.Response{
			Code:    int32(codes.OK),
			Message: codes.OK.String(),
			Payload: response,
		}, nil
	}
	for _, v := range s.endpoints {
		if ok := slices.ContainsFunc(v.sd.Methods, func(e grpc.MethodDesc) bool {
			return req.Method == e.MethodName
		}); !ok {
			continue
		}
		method := path.Join("/", v.sd.ServiceName, req.Method)
		reply, err := v.RoundTrip(ctx, &api.Request{Method: method, Payload: req.Payload, Timeout: req.Timeout})
		if err != nil {
			return &kubeapi.Response{
				Code:    int32(codes.Unavailable),
				Message: err.Error(),
			}, nil
		}
		return &kubeapi.Response{
			Code:    reply.Code,
			Message: reply.Message,
			Payload: reply.Payload,
		}, nil
	}

	return &kubeapi.Response{
		Code:    int32(codes.Unimplemented),
		Message: codes.Unimplemented.String(),
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
			var request kubeapi.Request
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
