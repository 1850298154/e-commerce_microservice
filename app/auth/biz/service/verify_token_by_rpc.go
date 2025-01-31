package service

import (
	"2501YTC/app/auth/biz/middlewares"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
	"errors"
	"time"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	j := middlewares.NewJWT()

	claims, err := j.ParseToken(req.Token)
	if err != nil {
		switch err {
		case middlewares.ErrTokenExpired:
			return nil, errors.New("token 已过期")
		case middlewares.ErrTokenMalformed:
			return nil, errors.New("token 格式错误")
		case middlewares.ErrTokenNotValidYet:
			return nil, errors.New("token 尚未激活")
		case middlewares.ErrTokenInvalid:
			return nil, errors.New("无效的 token")
		default:
			return nil, errors.New("token 验证失败")
		}
	}

	if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token 已过期")
	}

	return &auth.VerifyResp{
		Res:    true,
		UserId: int32(claims.UserId),
		Role:   claims.Role,
	}, nil
}
