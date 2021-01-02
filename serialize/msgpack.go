package serialize

import (
	msgpack "github.com/vmihailenco/msgpack/v5"
)

// MsgPackCoder implements ByteCoder interface
type MsgPackCoder struct{}

// Encode ..
func (e *MsgPackCoder) Encode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decode ..
func (e *MsgPackCoder) Decode(data []byte, vPtr interface{}) error {
	return msgpack.Unmarshal(data, vPtr)
}
