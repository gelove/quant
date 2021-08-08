package errs

import (
	"errors"
	"net/http"
)

// HTTP error code
const (
	StatusBadRequest       = http.StatusBadRequest
	StatusUnauthorized     = http.StatusUnauthorized
	StatusPaymentRequired  = http.StatusPaymentRequired
	StatusForbidden        = http.StatusForbidden
	StatusNotFound         = http.StatusNotFound
	StatusMethodNotAllowed = http.StatusMethodNotAllowed
	StatusNotAcceptable    = http.StatusNotAcceptable
)

// custom error code
const (
	Unknown = 1000 + iota
	DataIsNotFound
	DataIsNotEnough
)

var customInfo = map[int]string{
	DataIsNotFound:  "数据不存在",
	DataIsNotEnough: "刚上线的token不用统计, 数据太少",
}

func getInfo(code int) string {
	val, ok := customInfo[code]
	if !ok {
		return "未知错误"
	}
	return val
}

// customErrors 自定义错误
var customErrors = map[int]error{
	StatusBadRequest:       errors.New(http.StatusText(StatusBadRequest)),
	StatusUnauthorized:     errors.New(http.StatusText(StatusUnauthorized)),
	StatusPaymentRequired:  errors.New(http.StatusText(StatusPaymentRequired)),
	StatusForbidden:        errors.New(http.StatusText(StatusForbidden)),
	StatusNotFound:         errors.New(http.StatusText(StatusNotFound)),
	StatusMethodNotAllowed: errors.New(http.StatusText(StatusMethodNotAllowed)),
	StatusNotAcceptable:    errors.New(http.StatusText(StatusNotAcceptable)),
	DataIsNotFound:         errors.New(getInfo(DataIsNotFound)),
	DataIsNotEnough:        errors.New(getInfo(DataIsNotEnough)),
}

// Get Error
func Get(code int) error {
	val, ok := customErrors[code]
	if !ok {
		return errors.New("未知错误")
	}
	return val
}

// Check Error
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
