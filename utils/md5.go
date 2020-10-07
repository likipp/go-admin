package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/argon2"
)

var salt = []byte{0xc8, 0x28, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}

func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// 需要使用base64.StdEncoding.EncodeToString, 直接使用string会失败
func PasswordHash(password string) string {
	bytes := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	fmt.Println(base64.StdEncoding.EncodeToString(bytes), "string(bytes)", password, bytes)
	return base64.StdEncoding.EncodeToString(bytes)
}

func PasswordVerify(hash, password string) bool {
	return PasswordHash(password) == hash
}
