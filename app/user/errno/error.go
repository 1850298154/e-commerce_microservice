package errno

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// System Code
	ServiceErrCode = 10001
	ParamErrCode   = 10002

	// User ErrCode
	LoginErrCode            = 11001
	UserNotExistErrCode     = 11002
	UserAlreadyExistErrCode = 11003
	CreateUserErrCode       = 11004
	UpdateUserErrCode       = 11005
	DeleteUserErrCode       = 11006
)

func NewBizErr(err error, errCode int64, errMsg string) error {
	return errors.Wrap(err, fmt.Sprintf("err_code=%d, err_msg=%s", errCode, errMsg))
}

var (
	ServiceErr          = func(err error) error { return NewBizErr(err, ServiceErrCode, "服务没有重新启动") }
	ParamErr            = func(err error) error { return NewBizErr(err, ParamErrCode, "无效的参数请求") }
	LoginErr            = func(err error) error { return NewBizErr(err, LoginErrCode, "密码错误，登录失败") }
	UserNotExistErr     = func(err error) error { return NewBizErr(err, UserNotExistErrCode, "用户不存在") }
	UserAlreadyExistErr = func(err error) error { return NewBizErr(err, UserAlreadyExistErrCode, "用户已经存在") }
	CreateUserErr       = func(err error) error { return NewBizErr(err, CreateUserErrCode, "创建用户失败") }
	UpdateUserErr       = func(err error) error { return NewBizErr(err, UpdateUserErrCode, "更新用户失败") }
	DeleteUserErr       = func(err error) error { return NewBizErr(err, DeleteUserErrCode, "删除用户失败") }
)
