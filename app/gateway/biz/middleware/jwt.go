package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/cloudwego/hertz/pkg/app"
)

type CustomClaims struct {
	UserID uint32
	Role   uint32
	jwt.RegisteredClaims
}

type contextKey string

const Useridkey contextKey = "user_id"

var publicRoutes = map[string]struct{}{
	"/auth/token":     {},
	"/auth/verify":    {},
	"/auth/renew":     {},
	"/user/register":  {},
	"/user/login":     {},
	"/products":       {},
	"/product":        {},
	"/product/search": {},
}

func JwtAuthMiddleware(jwtSecret string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 公共路由处理
		path := string(c.Request.URI().Path())
		if _, ok := publicRoutes[path]; ok {
			c.Next(ctx)
			return
		}

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
		ctx = context.WithValue(ctx, Useridkey, claims.UserID)
		c.Next(ctx)
	}
}
