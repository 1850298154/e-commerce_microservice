package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserID uint32
	Role   uint32
	jwt.StandardClaims
}

func JwtAuthMiddleware(jwtSecret string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authStr := string(authHeader)
		tokenParts := strings.Split(authStr, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next(ctx)
	}
}
