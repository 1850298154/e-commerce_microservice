package service

import (
	"2501YTC/app/auth/errno"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/golang-jwt/jwt/v5"

	"2501YTC/app/auth/biz/dal/redis"

	"github.com/google/uuid"

	"2501YTC/app/auth/biz/middlewares"
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
	blacklistKey := fmt.Sprintf("jti_blacklist:%s", claims.RegisteredClaims.ID)
	exists, err := redis.RedisClient.Exists(s.ctx, blacklistKey).Result()
	if err != nil {
		return nil, err
	}
	if exists == 1 {
		err = errno.TokenRevokedErr(err)
		klog.Error(err)
		return nil, err
	}
	// 生成新的 JTI

	accessTokenRemaining := time.Until(claims.ExpiresAt.Time)
	refreshTokenRemaining := time.Until(claims.RegisteredClaims.ExpiresAt.Time)

	if accessTokenRemaining > 10*time.Minute && refreshTokenRemaining > 24*time.Hour {
		newClaims := *claims
		newAccessToken, err := j.CreateToken(newClaims)
		if err != nil {
			return nil, err
		}
		return &auth.RenewTokenResp{
			Token:        newAccessToken,
			RefreshToken: req.RefreshToken,
			ExpiresIn:    int64(accessTokenRemaining.Seconds()),
		}, nil
	}

	newJTI := uuid.New().String()

	// 生成新的 AccessToken
	newClaims := *claims
	newClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	newClaims.RegisteredClaims.ID = newJTI
	newAccessToken, err := j.CreateToken(newClaims)
	if err != nil {
		return nil, err
	}

	// 生成新的 RefreshToken
	newRefreshClaims := newClaims
	newRefreshClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)) // 7 天
	newRefreshToken, err := j.CreateToken(newRefreshClaims)
	if err != nil {
		return nil, err
	}
	fmt.Println("token生成完毕")
	// 旧refreshtoken加入黑名单
	if err := redis.RedisClient.Set(s.ctx, blacklistKey, "revoked", time.Until(claims.ExpiresAt.Time)).Err(); err != nil {
		return nil, err
	}
	fmt.Println("加入到黑名单")

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
