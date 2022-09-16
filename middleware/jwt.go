package middleware

import (
	"ChainServer/api/request"
	"errors"
	"net/http"

	"ChainServer/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		//modelToken := model.JwtBlacklist{
		//	Jwt: token,
		//}
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":   -2,
				"msg":    "未登录或非法访问，无权限",
				"result": nil,
			})
		}
		//if services.IsBlacklist(token, modelToken) {
		//	response.Result(response.ERROR, gin.H{
		//		"reload": true,
		//	}, "异地登陆或令牌失效", c)
		//	c.Abort()
		//	return
		//}
		j := NewJWT()
		//解析token
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == ErrTokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"code":   -3,
					"msg":    "授权已过期",
					"result": nil,
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code":   -4,
				"msg":    err.Error(),
				"result": nil,
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenExpired     error = errors.New("token is expired")
	ErrTokenNotValidYet error = errors.New("errFoo Token not active yet")
	ErrTokenMalformed   error = errors.New("that's not even a token")
	ErrTokenInvalid     error = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(config.Conf.Jwt.SigningKey),
	}
}

//创建Token
func (j *JWT) CreateToken(claims request.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//解析Token
func (j *JWT) ParseToken(tokenString string) (*request.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.UserClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.UserClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	}
	return nil, ErrTokenInvalid
}

//刷新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &request.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*request.UserClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", ErrTokenInvalid
}
