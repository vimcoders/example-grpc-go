package balance

import (
	"context"
	"example/generated/kubeapi"
	"log/slog"

	"example/encoding"

	"github.com/vimcoders/grpcx/metadata"

	"github.com/nats-io/nats.go"
)

func (s *Server) SubscriptionPush(ctx context.Context) (*nats.Subscription, error) {
	subj := kubeapi.PushService_ServiceDesc.ServiceName + ".>"
	return s.nc.Subscribe(subj, s.subscription)
}

func (s *Server) subscription(m *nats.Msg) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("nats subscription panic", "subject=", m.Subject, "err=", r)
		}
	}()
	var kv []string
	for k, v := range m.Header {
		if len(v) <= 0 {
			continue
		}
		kv = append(kv, k, v[0])
	}
	ctx := metadata.AppendToContext(context.Background(), kv...)
	codec := encoding.GetCodec(encoding.Name())
	desc := kubeapi.PushService_ServiceDesc
	for _, v := range desc.Methods {
		methodName := desc.ServiceName + "." + v.MethodName
		if methodName != m.Subject {
			continue
		}
		v.Handler(s, ctx, func(v any) error { return codec.Unmarshal(m.Data, v) }, nil)
	}
}

func (s *Server) PushChat(context.Context, *kubeapi.PushChatResponse) (*kubeapi.PushChatResponse, error) {
	return nil, nil
}
