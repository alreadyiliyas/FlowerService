package authcontext

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const CodecName = "json"

type jsonCodec struct{}

func (jsonCodec) Name() string {
	return CodecName
}

func (jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {
	encoding.RegisterCodec(jsonCodec{})
}

func Codec() encoding.Codec {
	return jsonCodec{}
}
