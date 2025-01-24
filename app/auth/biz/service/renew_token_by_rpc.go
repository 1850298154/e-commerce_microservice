package service

import (
	"2501YTC/app/auth/biz/middlewares"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
)

type RenewTokenByRPCService struct {
	ctx context.Context
} // NewRenewTokenByRPCService new RenewTokenByRPCService
func NewRenewTokenByRPCService(ctx context.Context) *RenewTokenByRPCService {
	return &RenewTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *RenewTokenByRPCService) Run(req *auth.RenewTokenReq) (resp *auth.RenewTokenResp, err error) {
	// todo

	jwtService := middlewares.NewJWT()

	newToken, err := jwtService.RefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &auth.RenewTokenResp{
		Token: newToken,
	}, nil
}
