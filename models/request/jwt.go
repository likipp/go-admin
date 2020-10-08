package request

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UUID       string
	ID         int
	Username   string
	NickName   string
	BufferTime int64
	jwt.StandardClaims
}
