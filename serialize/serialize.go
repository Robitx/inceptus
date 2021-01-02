// Package serialize provides unified interface for encoding and decoding data with several implementations (JSON, MsgPack, Gob, ..).
//
// Example usage:
//
//    var err error
//    var byteCoder serialize.ByteCoder
//    byteCoder = &serialize.JSONCoder{}
//    byteCoder = &serialize.MsgPackCoder{}
//
//    type MyStruct struct {
//    	A int
//    	B string
//    }
//
//    dataToSerialize := MyStruct{A: 100, B: "hey"}
//    encodedBytes, err := byteCoder.Encode(dataToSerialize)
//    if err != nil {
//    	fmt.Fprintf(os.Stderr, "serialization failed: %v\n", err)
//    	os.Exit(1)
//    }
//
//    decodedData := MyStruct{}
//    err = byteCoder.Decode(encodedBytes, &decodedData)
//    if err != nil {
//    	fmt.Fprintf(os.Stderr, "deserialization failed: %v\n", err)
//    	os.Exit(1)
//    }
//
//    fmt.Printf(
//    	"Original: %v\nSerialized: %v\nDecoded: %v\n",
//    	dataToSerialize, encodedBytes, decodedData)
//
// Will print out:
// Original: {100 hey}
// Serialized: [130 161 65 100 161 66 163 104 101 121]
// Decoded: {100 hey}
package serialize

// ByteCoder interface for easy switch between serialization formats (JSON, MsgPack, ..)
type ByteCoder interface {
	Encode(data interface{}) ([]byte, error)
	Decode(data []byte, vPtr interface{}) error
}
