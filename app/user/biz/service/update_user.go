package service

import (
	"context"

	"2501YTC/app/user/biz/dal/mysql"
	"2501YTC/app/user/biz/model"
	"2501YTC/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserService struct {
	ctx context.Context
} // NewUpdateUserService new UpdateUserService
func NewUpdateUserService(ctx context.Context) *UpdateUserService {
	return &UpdateUserService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserService) Run(req *user.UpdateUserReq) (resp *user.UpdateUserResp, err error) {
	// Finish your business logic.
	query := model.NewUserQuery(s.ctx, mysql.DB)
	newUser, err := query.GetUserById(req.GetUserId())
	if err != nil {
		klog.Error(err)
		return &user.UpdateUserResp{Success: false}, err
	}
	newUser.Email = req.Email
	pwdHashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		klog.Error(err)
		return &user.UpdateUserResp{Success: false}, err
	}
	newUser.PasswordHashed = string(pwdHashed)
	if err = query.UpdateUser(newUser); err != nil {
		klog.Error(err)
		return &user.UpdateUserResp{Success: false}, err
	}
	return &user.UpdateUserResp{Success: true}, nil
}
