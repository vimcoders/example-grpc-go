package example_test

import (
	"bufio"
	"crypto/tls"
	"example/generated/kubeapi"
	"log/slog"
	"math"
	"path"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func newTLSChannel() *channel {
	for range math.MaxInt8 {
		c, err := tls.Dial("tcp", "localhost:56888", &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		return &channel{
			Conn: c,
			br:   bufio.NewReader(c),
		}
	}
	panic("outofrange")
}

func BenchmarkTLSHello(b *testing.B) {
	ch := newTLSChannel()
	req := kubeapi.Request{
		Method:  path.Base(kubeapi.BalanceService_Hello_FullMethodName),
		Timeout: int64(time.Second),
	}
	s, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ch.Send(s); err != nil {
			continue
		}
		if _, err := ch.Recv(); err != nil {
			continue
		}
	}
}

func BenchmarkTLSChat(b *testing.B) {
	ch := newTLSChannel()
	req := kubeapi.Request{
		Method:  path.Base(kubeapi.ChatService_Chat_FullMethodName),
		Timeout: int64(time.Second),
	}
	s, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ch.Send(s); err != nil {
			continue
		}
		if _, err := ch.Recv(); err != nil {
			continue
		}
	}
}

func BenchmarkTLSLogin(b *testing.B) {
	ch := newTLSChannel()
	req := kubeapi.Request{
		Method:  path.Base(kubeapi.ProxyService_Login_FullMethodName),
		Timeout: int64(time.Second),
	}
	s, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ch.Send(s); err != nil {
			continue
		}
		if _, err := ch.Recv(); err != nil {
			continue
		}
	}
}

func BenchmarkTLSGetFriends(b *testing.B) {
	ch := newTLSChannel()
	req := kubeapi.Request{
		Method:  path.Base(kubeapi.SocialService_GetFriends_FullMethodName),
		Timeout: int64(time.Second),
	}
	s, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ch.Send(s); err != nil {
			continue
		}
		if _, err := ch.Recv(); err != nil {
			continue
		}
	}
}

func BenchmarkTLSGetActivity(b *testing.B) {
	ch := newTLSChannel()
	req := kubeapi.Request{
		Method:  path.Base(kubeapi.ActivityService_GetActivity_FullMethodName),
		Timeout: int64(time.Second),
	}
	s, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ch.Send(s); err != nil {
			continue
		}
		if _, err := ch.Recv(); err != nil {
			continue
		}
	}
}

func BenchmarkTLSGetMail(b *testing.B) {
	ch := newTLSChannel()
	req := kubeapi.Request{
		Method:  path.Base(kubeapi.MailService_GetMail_FullMethodName),
		Timeout: int64(time.Second),
	}
	s, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ch.Send(s); err != nil {
			continue
		}
		if _, err := ch.Recv(); err != nil {
			continue
		}
	}
}

func BenchmarkTLSDownload(b *testing.B) {
	ch := newTLSChannel()
	req := kubeapi.Request{
		Method:  path.Base(kubeapi.ItemService_Download_FullMethodName),
		Timeout: int64(time.Second),
	}
	s, err := proto.Marshal(&req)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := ch.Send(s); err != nil {
			continue
		}
		if _, err := ch.Recv(); err != nil {
			continue
		}
	}
}
