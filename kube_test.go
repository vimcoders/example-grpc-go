package kube_test

import (
	"bufio"
	"encoding/binary"
	"io"
	"kube/generated/kubeapi"
	"math"
	"net"
	"path"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	messageHeaderLength = 4
	messageLengthMax    = math.MaxUint16
)

type channel struct {
	net.Conn
	br *bufio.Reader
}

// NewChannel creates a new channel with the given net.Conn.
func newChannel() *channel {
	for {
		c, err := net.Dial("tcp", "localhost:26888")
		if err != nil {
			continue
		}
		return &channel{
			Conn: c,
			br:   bufio.NewReader(c),
		}
	}
}

// recv a message from the channel. The returned buffer contains the message.
//
// If a valid grpc status is returned, the message header
// returned will be valid and caller should send that along to
// the correct consumer. The bytes on the underlying channel
// will be discarded.
func (ch *channel) Recv() ([]byte, error) {
	var hrbuf [messageHeaderLength]byte // avoid alloc when reading header
	_, err := io.ReadFull(ch.br, hrbuf[:])
	if err != nil {
		return nil, err
	}

	h := binary.BigEndian.Uint32(hrbuf[:4])

	if h > uint32(messageLengthMax) {
		if _, err := ch.br.Discard(int(h)); err != nil {
			return nil, err
		}

		return nil, status.Error(codes.OutOfRange, codes.OutOfRange.String())
	}

	p := make([]byte, h)
	if _, err := io.ReadFull(ch.br, p); err != nil {
		return nil, err
	}

	return p, nil
}

// Send sends a message to the channel. The message is prefixed with a fixed-length header containing the length of the message and the stream ID.
func (ch *channel) Send(p []byte) error {
	if len(p) > messageLengthMax {
		return status.Error(codes.OutOfRange, codes.OutOfRange.String())
	}
	hwbuf := make([]byte, len(p)+messageHeaderLength)
	binary.BigEndian.PutUint32(hwbuf[:messageHeaderLength], uint32(len(p)))
	copy(hwbuf[messageHeaderLength:], p)

	_, err := ch.Write(hwbuf)
	if err != nil {
		return err
	}
	return nil
}

func BenchmarkHello(b *testing.B) {
	ch := newChannel()
	req := kubeapi.Request{Method: path.Base(kubeapi.HelloService_HelloEcho_FullMethodName)}
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

func BenchmarkKube(b *testing.B) {
	ch := newChannel()
	req := kubeapi.Request{Method: path.Base(kubeapi.KubeService_KubeEcho_FullMethodName)}
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
