package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId uint32
	Role   uint32
	jwt.RegisteredClaims
}

func (c *CustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.ExpiresAt, nil
}

func (c *CustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.IssuedAt, nil
}

func (c *CustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.NotBefore, nil
}
