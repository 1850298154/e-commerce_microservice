package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserId uint32
	Role   uint32
	jwt.StandardClaims
}
