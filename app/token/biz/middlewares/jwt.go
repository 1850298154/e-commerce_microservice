package middlewares

import (
	"errors"
	"time"

	"2501YTC/app/token/conf"

	models "2501YTC/app/token/biz/model"

	// "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 实现JWT服务
type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(conf.GetConf().JWT.SigningKey), // 可以设置过期时间

	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	// 生成token的前两段
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	// 使用签名，生成token
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析token包含的信息
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (any, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
		return j.CreateToken(*claims)
	}
	return "", ErrTokenInvalid
}
