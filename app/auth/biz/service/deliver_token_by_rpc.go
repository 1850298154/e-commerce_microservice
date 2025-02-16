package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"2501YTC/app/auth/biz/dal/redis"

	"github.com/google/uuid"

	"2501YTC/app/auth/biz/dal/mysql"
	"2501YTC/app/auth/biz/middlewares"
	models "2501YTC/app/auth/biz/model"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

type DeliverTokenByRPCService struct {
	ctx context.Context
} // NewDeliverTokenByRPCService new DeliverTokenByRPCService

func NewDeliverTokenByRPCService(ctx context.Context) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	if mysql.DB == nil {
		return nil, fmt.Errorf("mysql.DB is nil")
	}
	j := middlewares.NewJWT()
	jti := uuid.New().String()
	claims := models.CustomClaims{
		UserId: req.UserId,
		Role:   req.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    // 生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1小时过期
			Issuer:    "gomall",
			ID:        jti,
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}

	refreshClaims := claims
	refreshClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)) // 7天过期
	refreshToken, err := j.CreateToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	// 数据库保存
	tokenRecord := models.Token{
		UserID: req.UserId,
		Role:   req.Role,
	}
	tokenQuery := models.NewTokenQuery(s.ctx, mysql.DB)
	_, err = tokenQuery.Create(tokenRecord)
	if err != nil {
		return nil, err
	}
	// redis
	redisKey := fmt.Sprintf("user:%d:current_jti", req.UserId)
	if err := redis.RedisClient.Set(s.ctx, redisKey, jti, 7*24*time.Hour).Err(); err != nil {
		return nil, err
	}
	return &auth.DeliveryResp{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
