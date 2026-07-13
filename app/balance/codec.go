package balance

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// GetCodec returns a new instance of the proto codec.
type codec struct{}

// Name implements [encoding.Codec].
func (c *codec) Name() string {
	return "proto"
}

// Marshal implements [encoding.Codec].
func (c codec) Marshal(msg any) ([]byte, error) {
	switch v := msg.(type) {
	case proto.Message:
		return proto.Marshal(v)
	default:
		return nil, fmt.Errorf("ttrpc: cannot marshal unknown type: %T", msg)
	}
}

// Unmarshal implements [encoding.Codec].
func (c codec) Unmarshal(p []byte, msg any) error {
	switch v := msg.(type) {
	case proto.Message:
		return proto.Unmarshal(p, v)
	default:
		return fmt.Errorf("ttrpc: cannot unmarshal into unknown type: %T", msg)
	}
}
