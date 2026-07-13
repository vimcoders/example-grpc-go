package balance

import (
	"context"
	"kube/generated/kubeapi"
	"net"
	"sync"

	"github.com/vimcoders/grpcx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

type Option func(*Server)

type Server struct {
	encoding.Codec
	wg        sync.WaitGroup
	listener  net.Listener
	closed    context.CancelFunc
	endpoints []RoundTripper
}

func NewServer(opt ...Option) *Server {
	var s = Server{
		Codec: &codec{},
	}
	for i := range opt {
		opt[i](&s)
	}
	return &s
}

func (s *Server) ListenAndServe(ctx context.Context, addr string, opt ...Option) error {
	for i := range opt {
		opt[i](s)
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	cancelCtx, closed := context.WithCancel(ctx)
	s.closed = closed
	defer s.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		s.wg.Go(func() {
			session := Session{
				endpoints: s.endpoints,
				desc:      &kubeapi.HelloService_ServiceDesc,
				Codec:     &codec{},
			}
			_ = session.Handle(cancelCtx, conn)
		})
	}
}

// RegisterService registers a service and its implementation to the gRPC
// server. It is called from the IDL generated code. This must be called before
// invoking Serve. If ss is non-nil (for legacy code), its type is checked to
// ensure it implements sd.HandlerType.
func (s *Server) RegisterService(sd *grpc.ServiceDesc, endpoint string) error {
	cc, err := grpcx.Dial(endpoint)
	if err != nil {
		return err
	}
	s.endpoints = append(s.endpoints, RoundTripper{
		sd:                  sd,
		ClientConnInterface: cc,
	})
	return nil
}

func (s *Server) Close() error {
	if s.closed != nil {
		s.closed()
	}
	if s.listener != nil {
		s.listener.Close()
	}
	s.wg.Wait()
	return nil
}
