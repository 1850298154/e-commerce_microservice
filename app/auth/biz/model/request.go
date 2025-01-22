package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserId      uint
	AuthorityID uint
	jwt.StandardClaims
}
