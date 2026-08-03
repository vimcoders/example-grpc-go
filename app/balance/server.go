package balance

import (
	"context"
	"crypto/tls"
	"example/generated/kubeapi"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"example/encoding"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/vimcoders/grpcx"
	"google.golang.org/grpc"
)

type Option func(*Server)

type Options struct {
	encoding.Codec
	nc              *nats.Conn
	universalClient redis.UniversalClient
}

var defaultOptions = Options{
	Codec: encoding.GetCodec(encoding.Name()),
}

type Server struct {
	Options
	desc      *grpc.ServiceDesc
	wg        sync.WaitGroup
	listener  net.Listener
	closed    context.CancelFunc
	endpoints []RoundTripper
	kubeapi.UnimplementedPushServiceServer
}

func WithRedisService(e string) Option {
	if len(e) == 0 {
		return func(s *Server) {}
	}
	const defaultPoolSize = 8
	const defaultMaxRedirects = 3
	const defaultFailingTimeoutSeconds = 15
	opts := redis.UniversalOptions{
		Addrs: strings.Split(e, ","),
		// 每个节点的连接池大小
		PoolSize: defaultPoolSize,
		// 集群重定向最大次数（自动发现时可能需要）
		MaxRedirects: defaultMaxRedirects,
		// 节点失败标记时间（自动避开故障节点）
		FailingTimeoutSeconds: defaultFailingTimeoutSeconds,
		// 按延迟路由（可选，提升性能）
		RouteByLatency: true,
	}
	c := redis.NewUniversalClient(&opts)
	return func(s *Server) {
		if len(e) > 0 {
			s.universalClient = c
		}
	}
}

func WithNatsService(e string) (Option, error) {
	if len(e) == 0 {
		return func(s *Server) {}, nil
	}
	nc, err := nats.Connect(e)
	if err != nil {
		return nil, err
	}
	return func(s *Server) {
		s.nc = nc
	}, nil
}

func NewServer(opt ...Option) *Server {
	opts := defaultOptions
	var s = Server{
		Options: opts,
		desc:    &kubeapi.BalanceService_ServiceDesc,
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
				Codec:     s.Codec,
			}
			if err := session.Handle(cancelCtx, conn); err != nil {
				slog.Error("Handle", "error", err)
			}
		})
	}
}

func (s *Server) ListenAndServeTLS(ctx context.Context, addr string, opt ...Option) error {
	for i := range opt {
		opt[i](s)
	}
	certData, keyData, err := GenerateSignedCert()
	if err != nil {
		panic(err)
	}
	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		panic(err)
	}
	listener, err := tls.Listen("tcp", addr, &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})
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
				Codec:     s.Codec,
			}
			if err := session.Handle(cancelCtx, conn); err != nil {
				slog.Error("Handle", "error", err)
			}
		})
	}
}

// ServeHTTP 实现 http.Handler 接口
// 作用：提供 Admin 后台的 HTTP 接口服务，接收前端/运营后台的请求，并转发到对应的 Protobuf 方法
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/x-protobuf")
	// 如果是 OPTIONS 预检请求，直接返回 204，不处理业务
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	// 从请求头获取身份令牌（前端登录后的 token）
	// authority := r.Header.Get("Authorization")

	session := Session{
		endpoints: s.endpoints,
		desc:      s.desc,
		Codec:     s.Codec,
	}
	response, err := session.RoundTrip(context.Background(), &kubeapi.Request{
		Method:  path.Base(r.URL.Path),
		Payload: b,
		Timeout: int64(time.Second),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(response.Payload)
}

// RegisterService registers a service and its implementation to the gRPC
// server. It is called from the IDL generated code. This must be called before
// invoking Serve. If ss is non-nil (for legacy code), its type is checked to
// ensure it implements sd.HandlerType.
func (s *Server) RegisterService(sd grpc.ServiceDesc, endpoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cc, err := grpcx.DialContext(ctx, endpoint)
	if err != nil {
		return err
	}
	s.endpoints = append(s.endpoints, RoundTripper{
		ServiceDesc:         sd,
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
