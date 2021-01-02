package serialize

import (
	"bytes"
	"encoding/gob"
)

// GobCoder implements ByteCoder interface
type GobCoder struct{}

// Encode ..
func (e *GobCoder) Encode(data interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Decode ..
func (e *GobCoder) Decode(data []byte, vPtr interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(vPtr)
}
