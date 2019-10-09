package json

import (
	"github.com/json-iterator/go"
	"io"
)

var jsonConfig = jsoniter.ConfigCompatibleWithStandardLibrary

// Marshal json marshal
func Marshal(v interface{}) ([]byte, error) {
	return jsonConfig.Marshal(v)
}

// Unmarshal json unmarshal
func Unmarshal(data []byte, v interface{}) error {
	return jsonConfig.Unmarshal(data, v)
}

// NewDecoder json new decoder
func NewDecoder(reader io.Reader) *jsoniter.Decoder {
	return jsonConfig.NewDecoder(reader)
}

// NewEncoder json new encoder
func NewEncoder(writer io.Writer) *jsoniter.Encoder {
	return jsonConfig.NewEncoder(writer)
}
