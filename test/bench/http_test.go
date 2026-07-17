package example_test

import (
	"bytes"
	"example/encoding"
	"example/generated/kubeapi"
	"net/http"
	"testing"
	"time"
)

func BenchmarkHTTPHello(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.HelloRequest{Message: "HTTP Hello"})
	if err != nil {
		b.Error(err)
		return
	}
	client := http.Client{Timeout: time.Second}
	b.ResetTimer()
	for b.Loop() {
		if _, err := client.Post("http://127.0.0.1:16888"+kubeapi.BalanceService_Hello_FullMethodName, "application/x-protobuf", bytes.NewBuffer(req)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkHTTPChat(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.ChatRequest{Message: "HTTP Chat"})
	if err != nil {
		b.Error(err)
		return
	}
	client := http.Client{Timeout: time.Second}
	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Post("http://127.0.0.1:16888"+kubeapi.ChatService_Chat_FullMethodName, "application/x-protobuf", bytes.NewBuffer(req))
		if err != nil {
			b.Error(err)
		}
		_ = resp.Body.Close()
	}
}

func BenchmarkHTTPLogin(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.LoginRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	client := http.Client{Timeout: time.Second}
	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Post("http://127.0.0.1:16888"+kubeapi.ProxyService_Login_FullMethodName, "application/x-protobuf", bytes.NewBuffer(req))
		if err != nil {
			b.Error(err)
		}
		_ = resp.Body.Close()
	}
}

func BenchmarkHTTPGetFriends(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.GetFriendsRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	client := http.Client{Timeout: time.Second}
	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Post("http://127.0.0.1:16888"+kubeapi.SocialService_GetFriends_FullMethodName, "application/x-protobuf", bytes.NewBuffer(req))
		if err != nil {
			b.Error(err)
		}
		_ = resp.Body.Close()
	}
}

func BenchmarkHTTPGetActivity(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.GetActivityRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	client := http.Client{Timeout: time.Second}
	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Post("http://127.0.0.1:16888"+kubeapi.ActivityService_GetActivity_FullMethodName, "application/x-protobuf", bytes.NewBuffer(req))
		if err != nil {
			b.Error(err)
		}
		_ = resp.Body.Close()
	}
}

func BenchmarkHTTPGetMail(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.GetMailRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	client := http.Client{Timeout: time.Second}
	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Post("http://127.0.0.1:16888"+kubeapi.MailService_GetMail_FullMethodName, "application/x-protobuf", bytes.NewBuffer(req))
		if err != nil {
			b.Error(err)
		}
		_ = resp.Body.Close()
	}
}

func BenchmarkHTTPDownload(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.DownloadRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	client := http.Client{Timeout: time.Second}
	b.ResetTimer()
	for b.Loop() {
		resp, err := client.Post("http://127.0.0.1:16888"+kubeapi.ItemService_Download_FullMethodName, "application/x-protobuf", bytes.NewBuffer(req))
		if err != nil {
			b.Error(err)
		}
		_ = resp.Body.Close()
	}
}
