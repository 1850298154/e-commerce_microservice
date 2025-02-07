package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserId int32
	Role   int32
	jwt.StandardClaims
}
