package jwtauth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-admin/config"
	"go-admin/middleware"
	"go-admin/models/request"
	"go-admin/utils/response"
	"net/http"
	"strconv"
	"time"
)

type JWT struct {
	SigningKey []byte
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.Result(http.StatusNonAuthoritativeInfo, gin.H{
				"reload": true,
			}, "未登录或非法访问", false, c)
			c.Abort()
			return
		}
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == middleware.TokenExpired {
				response.Result(http.StatusExpectationFailed, gin.H{
					"reload": true,
				}, "授权已过期", false, c)
				c.Abort()
				return
			}
			response.Result(http.StatusExpectationFailed, gin.H{
				"reload": true,
			}, err.Error(), false, c)
			c.Abort()
			return
		}
		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
			claims.ExpiresAt = time.Now().Unix() + 60*60*24*7
			newToken, _ := j.CreateToken(*claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(config.AdminConfig.JWT.SigningKey),
	}
}

func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, middleware.TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, middleware.TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, middleware.TokenNotValidYet
			} else {
				return nil, middleware.TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, middleware.TokenInvalid

	} else {
		return nil, middleware.TokenInvalid

	}
}
