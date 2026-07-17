package example_test

import (
	"bytes"
	"example/encoding"
	"example/generated/kubeapi"
	"io"
	"net/http"
	"testing"
)

func post(path string, body io.Reader) error {
	resp, err := http.Post("http://localhost:16888"+path, "application/x-protobuf", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func BenchmarkHTTPHello(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.HelloRequest{Message: "HTTP Hello"})
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for b.Loop() {
		if err := post(kubeapi.BalanceService_Hello_FullMethodName, bytes.NewBuffer(req)); err != nil {
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
	b.ResetTimer()
	for b.Loop() {
		if err := post(kubeapi.ChatService_Chat_FullMethodName, bytes.NewBuffer(req)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkHTTPLogin(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	b.ResetTimer()
	for b.Loop() {
		req, err := codec.Marshal(&kubeapi.LoginRequest{})
		if err != nil {
			b.Error(err)
			return
		}
		if err := post(kubeapi.ProxyService_Login_FullMethodName, bytes.NewBuffer(req)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkHTTPGetFriends(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.GetFriendsRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for b.Loop() {
		if err := post(kubeapi.SocialService_GetFriends_FullMethodName, bytes.NewBuffer(req)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkHTTPGetActivity(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.GetActivityRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for b.Loop() {
		if err := post(kubeapi.ActivityService_GetActivity_FullMethodName, bytes.NewBuffer(req)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkHTTPGetMail(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.GetMailRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for b.Loop() {
		if err := post(kubeapi.MailService_GetMail_FullMethodName, bytes.NewBuffer(req)); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkHTTPDownload(b *testing.B) {
	codec := encoding.GetCodec(encoding.Name())
	req, err := codec.Marshal(&kubeapi.DownloadRequest{})
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for b.Loop() {
		if err := post(kubeapi.ItemService_Download_FullMethodName, bytes.NewBuffer(req)); err != nil {
			b.Error(err)
		}
	}
}
