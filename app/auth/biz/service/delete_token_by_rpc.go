package service

import (
	"2501YTC/app/auth/errno"
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"

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
		err = errno.TokenVoidErr(err)
		klog.Error(err)
		return nil, err
	}
	blacklistKey := fmt.Sprintf("jti_blacklist:%s", claims.RegisteredClaims.ID)
	expirationTime := time.Until(time.Unix(claims.ExpiresAt.Unix(), 0))
	if err := redis.RedisClient.Set(s.ctx, blacklistKey, "revoked", expirationTime).Err(); err != nil {
		err = errno.AddBlacklistTokenErr(err)
		klog.Error(err)
		return nil, err
	}

	return &auth.DeleteTokenResp{
		Res: true,
	}, nil
}
