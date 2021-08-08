package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"

	"github.com/pkg/errors"
)

// MD5 MD5哈希值
func MD5(b []byte) string {
	h := md5.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5String MD5哈希值
func MD5String(s string) string {
	return MD5([]byte(s))
}

// SHA1 SHA1哈希值
func SHA1(b []byte) string {
	h := sha1.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1String SHA1哈希值
func SHA1String(s string) string {
	return SHA1([]byte(s))
}

// SHA256 SHA256哈希值
func SHA256(b []byte) string {
	h := sha256.New()
	_, err := h.Write(b)
	if err != nil {
		panic(errors.WithStack(err))
	}
	hmac := []byte("NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j")
	return fmt.Sprintf("%x", h.Sum(hmac))
}

// SHA256String SHA256哈希值
func SHA256String(s string) string {
	return SHA256([]byte(s))
}

func SHA256MAC(s, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(s)
	// rs := mac.Sum(nil)
	rs := fmt.Sprintf("%x", mac.Sum(nil))
	return string(rs)
}

// func Hash(str string) (string, error) {
// 	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
// 	return string(hashed), err
// }

// func IsSame(str string, hashed string) bool {
// 	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str)) == nil
// }
