package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserId    uint
	Role      string
	Authority string
	jwt.StandardClaims
}
