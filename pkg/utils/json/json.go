package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// 定义JSON操作
var (
	_json         = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = _json.Marshal
	Unmarshal     = _json.Unmarshal
	MarshalIndent = _json.MarshalIndent
	NewDecoder    = _json.NewDecoder
	NewEncoder    = _json.NewEncoder
)

// MarshalToString JSON编码为字符串
func MarshalToString(data interface{}) string {
	s, err := Marshal(data)
	if err != nil {
		return ""
	}
	return string(s)
}

// MustTransform MustTransform
func MustTransform(data, out interface{}) {
	bytes := MustEncode(data)
	MustDecode(bytes, out)
}

// MustDecode MustDecode
func MustDecode(data []byte, out interface{}) {
	err := Unmarshal(data, out)
	if err != nil {
		// logging.Errorf("MustDecode data => %#v", string(data))
		panic(errors.Wrap(err, "MustDecode Error!"))
	}
}

// MustEncode MustEncode
func MustEncode(data interface{}) []byte {
	bytes, err := Marshal(data)
	if err != nil {
		// logging.Errorf("MustEncode data => %#v", data)
		panic(errors.Wrap(err, "MustEncode Error!"))
	}
	return bytes
}

// MustToString MustToString
func MustToString(data interface{}) string {
	bytes := MustEncode(data)
	return string(bytes)
}
