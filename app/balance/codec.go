package balance

import (
	"encoding/json"

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
		return json.Marshal(msg)
	}
}

// Unmarshal implements [encoding.Codec].
func (c codec) Unmarshal(p []byte, msg any) error {
	switch v := msg.(type) {
	case proto.Message:
		return proto.Unmarshal(p, v)
	default:
		return json.Unmarshal(p, v)
	}
}
