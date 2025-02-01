package service

import (
	"context"

	"2501YTC/app/user/biz/dal/mysql"
	"2501YTC/app/user/errno"

	"github.com/cloudwego/kitex/pkg/klog"

	"2501YTC/app/user/biz/model"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"2501YTC/rpc_gen/kitex_gen/user"
)

type RegisterService struct {
	ctx context.Context
} // NewRegisterService new RegisterService
func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

// Run create note info
func (s *RegisterService) Run(req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// Finish your business logic.
	if req.Password != req.ConfirmPassword {
		err = errors.New("两次输入的密码不一致")
		klog.Error(err)
		return nil, err
	}
	query := model.NewUserQuery(s.ctx, mysql.DB)
	u, _ := query.GetUserByEmail(req.GetEmail())
	if u != nil {
		err = errno.UserAlreadyExistErr(err)
		klog.Error(err)
		return nil, err
	}
	pwdHashed, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		err = errors.Wrap(err, "密码加密时产生错误")
		klog.Error(err)
		return nil, err
	}
	userID, err := query.CreateUser(&model.User{Email: req.Email, PasswordHashed: string(pwdHashed), Role: model.UserRole})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return &user.RegisterResp{
		UserId: userID,
		Role:   int32(model.UserRole),
	}, nil
}
