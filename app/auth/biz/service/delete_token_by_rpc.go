package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"2501YTC/app/auth/biz/dal/redis"
	"2501YTC/app/auth/biz/middlewares"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
)

type DeleteTokenByRPCService struct {
	ctx context.Context
} // NewDeleteTokenByRPCService new DeleteTokenByRPCService
func NewDeleteTokenByRPCService(ctx context.Context) *DeleteTokenByRPCService {
	return &DeleteTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeleteTokenByRPCService) Run(req *auth.DeleteTokenReq) (resp *auth.DeleteTokenResp, err error) {
	// Finish your business logic.
	j := middlewares.NewJWT()
	claims, err := j.ParseToken(req.Token)
	if err != nil {
		return nil, errors.New("无效的Token")
	}
	blacklistKey := fmt.Sprintf("jti_blacklist:%s", claims.StandardClaims.Id)
	if err := redis.RedisClient.Set(s.ctx, blacklistKey, "revoked", time.Until(time.Unix(claims.StandardClaims.ExpiresAt, 0))).Err(); err != nil {
		return nil, err
	}

	return &auth.DeleteTokenResp{
		Res: true,
	}, nil
}
