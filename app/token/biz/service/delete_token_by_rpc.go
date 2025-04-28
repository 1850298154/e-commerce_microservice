package service

import (
	"context"
	"fmt"
	"time"

	"2501YTC/app/token/errno"

	"github.com/cloudwego/kitex/pkg/klog"

	"2501YTC/app/token/biz/dal/redis"
	"2501YTC/app/token/biz/middlewares"
	token "2501YTC/rpc_gen/kitex_gen/token"
)

type DeleteTokenByRPCService struct {
	ctx context.Context
} // NewDeleteTokenByRPCService new DeleteTokenByRPCService
func NewDeleteTokenByRPCService(ctx context.Context) *DeleteTokenByRPCService {
	return &DeleteTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeleteTokenByRPCService) Run(req *token.DeleteTokenReq) (resp *token.DeleteTokenResp, err error) {
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

	return &token.DeleteTokenResp{
		Res: true,
	}, nil
}
