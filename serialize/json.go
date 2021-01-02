package serialize

import (
	"encoding/json"
	"strings"
)

// JSONCoder implements ByteCoder interface
type JSONCoder struct{}

// Encode ..
func (e *JSONCoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode ..
func (e *JSONCoder) Decode(data []byte, vPtr interface{}) (err error) {
	switch arg := vPtr.(type) {
	case *string:
		// If they want a string and it is a JSON string, strip quotes
		// This allows someone to send a struct but receive as a plain string
		// This cast should be efficient for Go 1.3 and beyond.
		str := string(data)
		if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
			*arg = str[1 : len(str)-1]
		} else {
			*arg = str
		}
	case *[]byte:
		*arg = data
	default:
		err = json.Unmarshal(data, arg)
	}
	return
}
