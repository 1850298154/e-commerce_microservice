package service

import (
	"2501YTC/app/auth/biz/middlewares"
	models "2501YTC/app/auth/biz/model"
	"2501YTC/app/order/biz/dal/mysql"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"errors"
	"time"
)

type RenewTokenByRPCService struct {
	ctx context.Context
} // NewRenewTokenByRPCService new RenewTokenByRPCService
func NewRenewTokenByRPCService(ctx context.Context) *RenewTokenByRPCService {
	return &RenewTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *RenewTokenByRPCService) Run(req *auth.RenewTokenReq) (resp *auth.RenewTokenResp, err error) {

	jwtService := middlewares.NewJWT()
	// 解析旧 RefreshToken
	claims, err := jwtService.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("无效的refreshToken")
	}
	// 生成新的 AccessToken
	newClaims := *claims
	newClaims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix() // 1 小时
	newAccessToken, err := jwtService.CreateToken(newClaims)
	if err != nil {
		return nil, err
	}
	// 生成新的 RefreshToken
	newRefreshClaims := *claims
	newRefreshClaims.StandardClaims.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix() // 7 天
	newRefreshToken, err := jwtService.CreateToken(newRefreshClaims)
	if err != nil {
		return nil, err
	}
	// 刷新数据库
	tokenQuery := models.NewTokenQuery(s.ctx, mysql.DB)
	tokenRecord, err := tokenQuery.GetByUserID(claims.UserId)
	if err != nil {
		return nil, err
	}
	tokenRecord.Token = newAccessToken
	tokenRecord.RefreshToken = newRefreshToken
	tokenRecord.AccessExpires = time.Now().Add(1 * time.Hour).Unix()       // 30天
	tokenRecord.RefreshExpires = time.Now().Add(7 * 24 * time.Hour).Unix() // 60天
	_, err = tokenQuery.Update(claims.UserId, tokenRecord)
	if err != nil {
		return nil, err
	}

	return &auth.RenewTokenResp{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    1 * 60 * 60, // 1小时
	}, nil
}
