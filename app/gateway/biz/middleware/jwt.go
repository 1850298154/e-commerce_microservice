package middleware

import (
	"2501YTC/app/gateway/biz/dal/mysql"
	"2501YTC/app/gateway/biz/service"
	conf2 "2501YTC/app/gateway/conf"
	"2501YTC/app/gateway/hertz_gen/gateway/auth"
	"2501YTC/app/user/biz/model"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudwego/kitex/pkg/klog"

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

func JwtAuthMiddleware(jwtSecret string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 公共路由处理
		conf := conf2.GetConf()
		publicRoutes := make(map[string]struct{})
		for _, route := range conf.Security.PublicRoutes {
			publicRoutes[route] = struct{}{}
		}
		path := string(c.Request.URI().Path())
		if _, ok := publicRoutes[path]; ok {
			c.Next(ctx)
			return
		}

		tokenHeader := c.Request.Header.Get("Authorization")
		refreshTokenHeader := c.Request.Header.Get("X-Refresh-Token")
		if tokenHeader == "" || refreshTokenHeader == "" {
			fmt.Println("1")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authStr := string(tokenHeader)
		authRefreshStr := string(refreshTokenHeader)
		if !strings.HasPrefix(authStr, "Bearer ") || !strings.HasPrefix(authRefreshStr, "Bearer ") {
			fmt.Println("2")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var req auth.VerifyTokenReq
		req.Token = authStr
		req.RefreshToken = authRefreshStr
		var renewResp auth.RenewTokenResp
		_, err := service.NewVerifyTokenByRPCService(ctx, c).Run(&req)
		if err != nil {
			var req auth.RenewTokenReq
			req.RefreshToken = authRefreshStr
			tempResp, err := service.NewRenewTokenByRPCService(ctx, c).Run(&req)
			if err != nil {
				fmt.Println("8")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			renewResp = *tempResp
			authStr = renewResp.Token
			authRefreshStr = renewResp.RefreshToken
		}
		authStr = authStr[len("Bearer "):]
		authRefreshStr = authRefreshStr[len("Bearer "):]

		token, err := jwt.ParseWithClaims(authStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})
		fmt.Println(token)
		if err != nil || !token.Valid {
			fmt.Println("3")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			fmt.Println("4")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID := claims.UserID
		// 检查黑名单
		query := model.NewUserQuery(ctx, mysql.DB)
		u, err := query.GetUserById(userID)
		fmt.Println(u)
		if err != nil {
			klog.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// if u.IsBanned {
		//	c.JSON(http.StatusForbidden, map[string]string{
		//		"error":   "UserBanned",
		//		"message": "该用户已被封禁，请联系管理员。",
		//	})
		//	return
		// }

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		ctx = context.WithValue(ctx, Useridkey, claims.UserID)
		c.Next(ctx)
	}
}
