package encoding

import (
	"encoding/json"
	"errors"

	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

type Codec encoding.Codec

func Name() string {
	return defaultCodec.Name()
}

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
		return nil, errors.New("protoc not supported")
	}
}

// Unmarshal implements [encoding.Codec].
func (c codec) Unmarshal(p []byte, msg any) error {
	switch v := msg.(type) {
	case proto.Message:
		return proto.Unmarshal(p, v)
	default:
		return errors.New("protoc not supported")
	}
}

var defaultJsonc jsonc
var defaultCodec codec

func GetCodec(codecName string) encoding.Codec {
	switch codecName {
	case defaultJsonc.Name():
		return &defaultJsonc
	case defaultCodec.Name():
		return &defaultCodec
	}

	return &defaultCodec
}

type jsonc struct{}

// Name implements [encoding.Codec].
func (c *jsonc) Name() string {
	return "json"
}

// Marshal implements [encoding.Codec].
func (c jsonc) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal implements [encoding.Codec].
func (c jsonc) Unmarshal(p []byte, v any) error {
	return json.Unmarshal(p, v)
}
