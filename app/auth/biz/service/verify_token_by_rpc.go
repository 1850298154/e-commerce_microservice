package service

import (
	"2501YTC/app/auth/biz/middlewares"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	jwtService := middlewares.NewJWT()
	//todu
	claims, err := jwtService.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	// Return the result
	return &auth.VerifyResp{
		Res: true,
	}, nil
}
