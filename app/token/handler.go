package main

import (
	"context"

	"2501YTC/app/token/biz/service"
	token "2501YTC/rpc_gen/kitex_gen/token"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *token.DeliverTokenReq) (resp *token.DeliveryResp, err error) {
	resp, err = service.NewDeliverTokenByRPCService(ctx).Run(req)

	return resp, err
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *token.VerifyTokenReq) (resp *token.VerifyResp, err error) {
	resp, err = service.NewVerifyTokenByRPCService(ctx).Run(req)

	return resp, err
}

// RenewTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RenewTokenByRPC(ctx context.Context, req *token.RenewTokenReq) (resp *token.RenewTokenResp, err error) {
	resp, err = service.NewRenewTokenByRPCService(ctx).Run(req)

	return resp, err
}

// DeleteTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeleteTokenByRPC(ctx context.Context, req *token.DeleteTokenReq) (resp *token.DeleteTokenResp, err error) {
	resp, err = service.NewDeleteTokenByRPCService(ctx).Run(req)

	return resp, err
}
