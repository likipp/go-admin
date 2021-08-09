package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-admin/config"
	"go-admin/init/cookies"
	"go-admin/models"
	"go-admin/utils/response"
	"net/http"
	"strconv"
	"time"
)

var (
	TokenExpired     error = errors.New("token is expired")
	TokenNotValidYet error = errors.New("token not active yet")
	TokenMalformed   error = errors.New("that's not even a token")
	TokenInvalid     error = errors.New("couldn't handle this token")
)

type JWT struct {
	SigningKey []byte
}

//type CustomClaims struct {
//	UUID       string
//	ID         int
//	Username   string
//	NickName   string
//	RoleName   []models.SysRole
//	BufferTime int64
//	jwt.StandardClaims
//}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := cookies.GetSession(c)
		if err != nil {
			c.Abort()
			return
		}
		if session.Options.MaxAge < 0 {
			response.FailWithMessage("session已过期, 请重新登录.", c)
			c.Abort()
			return
		}
		token, ok := session.Values["token"].(string)
		if !ok {
			response.Result(http.StatusExpectationFailed, gin.H{
				"reload": true,
			}, "未登录或token已过期", 2, false, c)
			c.Abort()
			return
		}
		//if token == "" {
		//	response.Result(http.StatusNonAuthoritativeInfo, gin.H{
		//		"reload": true,
		//	}, "未登录或非法访问", 2, false, c)
		//	c.Abort()
		//	return
		//}
		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				response.Result(http.StatusExpectationFailed, gin.H{
					"reload": true,
				}, "授权已过期", 2, false, c)
				c.Abort()
				return
			}
			response.Result(http.StatusExpectationFailed, gin.H{
				"reload": true,
			}, err.Error(), 2, false, c)
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

func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}
}

//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", nil
//	}
//	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}
