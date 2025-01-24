package service

import (
	"context"

	"2501YTC/app/user/biz/dal/mysql"
	"2501YTC/app/user/biz/model"
	"2501YTC/app/user/errno"
	"2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

type UpdateUserRoleService struct {
	ctx context.Context
} // NewUpdateUserRoleService new UpdateUserRoleService
func NewUpdateUserRoleService(ctx context.Context) *UpdateUserRoleService {
	return &UpdateUserRoleService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserRoleService) Run(req *user.UpdateUserRoleReq) (resp *user.UpdateUserRoleResp, err error) {
	// Finish your business logic.
	query := model.NewUserQuery(s.ctx, mysql.DB)
	u, err := query.GetUserById(req.GetUserId())
	if err != nil {
		err = errno.UserNotExistErr(err)
		klog.Error(err)
		return &user.UpdateUserRoleResp{Success: false}, err
	}
	u.Role = model.Role(req.GetRole())
	if err = query.UpdateUser(u); err != nil {
		err = errno.UpdateUserErr(err)
		klog.Error(err)
		return &user.UpdateUserRoleResp{Success: false}, err
	}
	return &user.UpdateUserRoleResp{Success: true}, nil
}
