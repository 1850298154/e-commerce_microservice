package service

import (
	"context"

	"2501YTC/app/user/biz/dal/mysql"
	"2501YTC/app/user/biz/model"

	"github.com/cloudwego/kitex/pkg/klog"

	"2501YTC/rpc_gen/kitex_gen/user"
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
	query := model.NewUserQuery(s.ctx, mysql.DB)
	if err := query.DeleteUser(req.GetUserId()); err != nil {
		klog.Error(err)
		return &user.DeleteUserResp{Success: false}, err
	}
	return &user.DeleteUserResp{Success: true}, nil
}
