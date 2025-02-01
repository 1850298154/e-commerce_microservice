package service

import (
	"context"

	"2501YTC/app/user/biz/dal/mysql"
	"2501YTC/app/user/biz/model"
	"2501YTC/app/user/errno"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"

	"2501YTC/rpc_gen/kitex_gen/user"
)

type LoginService struct {
	ctx context.Context
} // NewLoginService new LoginService
func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

// Run create note info
func (s *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	// Finish your business logic.
	query := model.NewUserQuery(s.ctx, mysql.DB)
	klog.Infof("登录信息: %v", req)
	u, err := query.GetUserByEmail(req.GetEmail())
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHashed), []byte(req.Password)); err != nil {
		err = errno.LoginErr(err)
		klog.Error(err)
		return nil, err
	}
	return &user.LoginResp{
		UserId: int32(u.ID),
		Role:   int32(u.Role),
	}, nil
}
