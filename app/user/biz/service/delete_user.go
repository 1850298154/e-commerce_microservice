package service

import (
	user "2501YTC/rpc_gen/kitex_gen/user"
	"context"
)

type DeleteUserService struct {
	ctx context.Context
} // NewDeleteUserService new DeleteUserService
func NewDeleteUserService(ctx context.Context) *DeleteUserService {
	return &DeleteUserService{ctx: ctx}
}

// Run create note info
func (s *DeleteUserService) Run(req *user.DeleteUserReq) (resp *user.DeleteUserResp, err error) {
	// Finish your business logic.

	return
}
