package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UUID       string `json:"uuid"`
	ID         uint
	UserName   string
	NickName   string
	Roles      []SysRole
	BufferTime int64
	jwt.StandardClaims
}
