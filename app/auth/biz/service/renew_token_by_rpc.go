package service

import (
	"2501YTC/app/auth/biz/dal/redis"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"2501YTC/app/auth/biz/middlewares"
	models "2501YTC/app/auth/biz/model"
	"2501YTC/app/order/biz/dal/mysql"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

type RenewTokenByRPCService struct {
	ctx context.Context
} // NewRenewTokenByRPCService new RenewTokenByRPCService
func NewRenewTokenByRPCService(ctx context.Context) *RenewTokenByRPCService {
	return &RenewTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *RenewTokenByRPCService) Run(req *auth.RenewTokenReq) (resp *auth.RenewTokenResp, err error) {
	j := middlewares.NewJWT()
	// 解析旧refreshToken
	claims, err := j.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("无效的refreshToken")
	}

	// 检查旧refreshToken是否在黑名单中
	blacklistKey := fmt.Sprintf("jti_blacklist:%s", claims.StandardClaims.Id)
	exists, err := redis.RedisClient.Exists(s.ctx, blacklistKey).Result()
	if err != nil {
		return nil, err
	}
	if exists == 1 {
		return nil, errors.New("refreshToken 已被撤销")
	}

	// 生成新的 JTI
	newJTI := uuid.New().String()

	// 生成新的 AccessToken
	newClaims := *claims
	newClaims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix() // 1 小时
	newClaims.StandardClaims.Id = newJTI
	newAccessToken, err := j.CreateToken(newClaims)
	if err != nil {
		return nil, err
	}
	// 生成新的 RefreshToken
	newRefreshClaims := newClaims
	newRefreshClaims.StandardClaims.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix() // 7 天
	newRefreshToken, err := j.CreateToken(newRefreshClaims)
	if err != nil {
		return nil, err
	}

	// 旧refreshtoken加入黑名单
	if err := redis.RedisClient.Set(s.ctx, blacklistKey, "revoked", time.Until(time.Unix(claims.StandardClaims.ExpiresAt, 0))).Err(); err != nil {
		return nil, err
	}

	// 刷新数据库
	tokenQuery := models.NewTokenQuery(s.ctx, mysql.DB)
	tokenRecord, err := tokenQuery.GetByUserID(claims.UserId)
	if err != nil {
		return nil, err
	}
	_, err = tokenQuery.Update(claims.UserId, tokenRecord)
	if err != nil {
		return nil, err
	}
	// 更新redis
	redisKey := fmt.Sprintf("user:%d:current_jti", claims.UserId)
	if err := redis.RedisClient.Set(s.ctx, redisKey, newJTI, 7*24*time.Hour).Err(); err != nil {
		return nil, err
	}

	return &auth.RenewTokenResp{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    1 * 60 * 60, // 1小时
	}, nil
}
