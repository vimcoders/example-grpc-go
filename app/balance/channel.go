package balance

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"net"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	messageHeaderLength = 4
	messageLengthMax    = math.MaxUint16
)

var buffers sync.Pool

// channel is a wrapper around a net.Conn that provides methods for sending and receiving messages with a fixed-length header.
type channel struct {
	net.Conn
	br *bufio.Reader
}

// NewChannel creates a new channel with the given net.Conn.
func newChannel(conn net.Conn) *channel {
	return &channel{
		Conn: conn,
		br:   bufio.NewReader(conn),
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

	var p []byte
	if h > 0 {
		p = ch.getmbuf(int(h))
		if _, err := io.ReadFull(ch.br, p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// Send sends a message to the channel. The message is prefixed with a fixed-length header containing the length of the message and the stream ID.
func (ch *channel) Send(p []byte) error {
	if len(p) > messageLengthMax {
		return status.Error(codes.OutOfRange, codes.OutOfRange.String())
	}
	hwbuf := ch.getmbuf(messageHeaderLength + len(p))
	defer ch.putmbuf(hwbuf)
	binary.BigEndian.PutUint32(hwbuf[:4], uint32(len(p)))
	copy(hwbuf[messageHeaderLength:], p)

	_, err := ch.Write(hwbuf)
	if err != nil {
		return err
	}
	return nil
}

// getmbuf returns a buffer from the pool. The buffer is guaranteed to be at least size bytes long.
func (ch *channel) getmbuf(size int) []byte {
	// we can't use the standard New method on pool because we want to allocate
	// based on size.
	b, ok := buffers.Get().(*[]byte)
	if !ok || cap(*b) < size {
		// TODO(stevvooe): It may be better to allocate these in fixed length
		// buckets to reduce fragmentation but its not clear that would help
		// with performance. An ilogb approach or similar would work well.
		return make([]byte, size)
	}
	*b = (*b)[:size]
	return *b
}

// putmbuf returns a buffer to the pool. The buffer must have been allocated by getmbuf.
func (ch *channel) putmbuf(p []byte) {
	buffers.Put(&p)
}
