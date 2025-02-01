package main

import (
	"2501YTC/app/auth/biz/service"
	auth "2501YTC/rpc_gen/kitex_gen/auth"
	"context"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	resp, err = service.NewDeliverTokenByRPCService(ctx).Run(req)

	return resp, err
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	resp, err = service.NewVerifyTokenByRPCService(ctx).Run(req)

	return resp, err
}

// RenewTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RenewTokenByRPC(ctx context.Context, req *auth.RenewTokenReq) (resp *auth.RenewTokenResp, err error) {
	resp, err = service.NewRenewTokenByRPCService(ctx).Run(req)

	return resp, err
}
