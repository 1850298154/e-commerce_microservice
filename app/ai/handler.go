package main

import (
	"2501YTC/app/ai/biz/service"
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// SearchforOrder implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) SearchforOrder(ctx context.Context, req *ai.SearchforOrderReq) (resp *ai.SearchforOrderResp, err error) {
	resp, err = service.NewSearchforOrderService(ctx).Run(req)

	return resp, err
}
