package service

import (
	"context"
	"errors"

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
	if req.Email == "" {
		err = errors.New("邮箱不能为空！")
		klog.Error(err)
		return
	}
	if req.Password == "" {
		err = errors.New("密码不能为空！")
		klog.Error(err)
		return
	}
	query := model.NewUserQuery(s.ctx, mysql.DB)
	klog.Infof("登录信息: %v", req)
	u, err := query.GetUserByEmail(req.GetEmail())
	if err != nil {
		klog.Error(err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHashed), []byte(req.Password)); err != nil {
		err = errno.LoginErr(err)
	}
	if u.IsBanned {
		err := errors.New("user was banned")
		err = errno.UserBannedErr(err)
		klog.Error(err)
		return
	}
	return &user.LoginResp{
		UserId: uint32(u.ID),
		Role:   uint32(u.Role),
	}, nil

}
