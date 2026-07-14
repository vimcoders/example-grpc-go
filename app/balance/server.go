package balance

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"
	"net"
	"net/http"
	"path"
	"slices"
	"sync"
	"time"

	"github.com/vimcoders/grpcx"
	"github.com/vimcoders/grpcx/generated/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
)

type Option func(*Server)

type Server struct {
	encoding.Codec
	desc        *grpc.ServiceDesc
	wg          sync.WaitGroup
	listener    net.Listener
	closed      context.CancelFunc
	endpoints   []RoundTripper
	interceptor grpc.UnaryServerInterceptor
}

func NewServer(opt ...Option) *Server {
	var s = Server{
		Codec: &codec{},
		desc:  &kubeapi.HelloService_ServiceDesc,
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
				desc:      s.desc,
				Codec:     &codec{},
			}
			_ = session.Handle(cancelCtx, conn)
		})
	}
}

// ServeHTTP 实现 http.Handler 接口
// 作用：提供 Admin 后台的 HTTP 接口服务，接收前端/运营后台的请求，并转发到对应的 Protobuf 方法
func (x *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 统一延迟处理：捕获 panic 崩溃 + 确保请求体关闭
	defer func() {
		// 捕获 HTTP 处理过程中的 panic，防止整个服务崩溃
		if e := recover(); e != nil {
			slog.Error("ServeHTTP", "recover", e)
		}
		// 确保请求体关闭，避免连接泄漏
		if err := r.Body.Close(); err != nil {
			slog.Error("ServeHTTP close", "err", err.Error())
		}
	}()
	// 允许所有域跨域（生产环境可改为具体域名）
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// 设置响应内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")
	// 如果是 OPTIONS 预检请求，直接返回 204，不处理业务
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// 从请求头获取身份令牌（前端登录后的 token）
	authority := r.Header.Get("Authorization")
	slog.Debug("ServeHTTP", "path", r.URL.Path, "authority", authority)
	response, err := x.RoundTrip(context.Background(), &api.Request{Method: path.Base(r.URL.Path)})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	reply, err := x.Codec.Marshal(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(reply)
}

func (s *Server) RoundTrip(ctx context.Context, req *api.Request) (*api.Response, error) {
	if req.Method == "" {
		return &api.Response{
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
	for _, v := range s.endpoints {
		if ok := slices.ContainsFunc(v.sd.Methods, func(e grpc.MethodDesc) bool {
			return req.Method == e.MethodName
		}); !ok {
			continue
		}
		method := path.Join("/", v.sd.ServiceName, req.Method)
		reply, err := v.RoundTrip(ctx, &api.Request{Method: method, Payload: req.Payload})
		slog.Info("RoundTrip", "reply", reply, "err", err)
		if err != nil {
			return &api.Response{
				Code:    int32(codes.Unavailable),
				Message: err.Error(),
			}, nil
		}
		return &api.Response{
			Code:    reply.Code,
			Message: reply.Message,
			Payload: reply.Payload,
		}, nil
	}

	return &api.Response{
		Code:    int32(codes.Unimplemented),
		Message: codes.Unimplemented.String(),
	}, nil
}

func (s *Server) HelloEcho(ctx context.Context, req *kubeapi.HelloRequest) (*kubeapi.HelloResponse, error) {
	slog.Info("echo", "hello", req)
	return &kubeapi.HelloResponse{Message: req.Message}, nil
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
