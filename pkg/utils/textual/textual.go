package textual

import (
	cRand "crypto/rand"
	"encoding/hex"
	"io"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

const (
	letterBytes   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString 获取指定长度随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	src := rand.NewSource(time.Now().UnixNano())
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// RandomKey 获取指定长度随机字符串
func RandomKey(length int) string {
	k := make([]byte, length)
	if _, err := io.ReadFull(cRand.Reader, k); err != nil {
		panic(err)
	}
	s := hex.EncodeToString(k)
	return s[:length]
}

// ArrayShift shift an element off the beginning of slice
func ArrayShift(s *[]string) string {
	if len(*s) == 0 {
		return ""
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// ArrayPop pop an element off the last of slice
func ArrayPop(s *[]string) string {
	l := len(*s)
	if l == 0 {
		return ""
	}
	f := (*s)[l-1]
	*s = (*s)[:l-1]
	return f
}

// InArray 判断是否在数组切片中
func InArray(needle string, list []string) bool {
	for _, v := range list {
		if strings.Trim(v, " ") == needle {
			return true
		}
	}
	return false
}

// TrimSpacing 去除元素前后的空格
func TrimSpacing(list []string) {
	for index := range list {
		list[index] = strings.Trim(list[index], " ")
	}
	return
}

// FilterSpacing 过滤空格元素
func FilterSpacing(list []string) []string {
	ret := make([]string, 0, len(list))
	for _, val := range list {
		if val != " " {
			ret = append(ret, val)
		}
	}
	return ret
}
