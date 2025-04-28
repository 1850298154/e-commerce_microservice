package errno

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// System Code
	ServiceErrCode = 10001
	ParamErrCode   = 10002

	// Auth ErrCode
	TokenVoidErrCode         = 12001
	TokenExpiredErrCode      = 12002
	DeleteTokenErrCode       = 12003
	CreateTokenErrCode       = 12004
	TokenRevokedErrCode      = 12005
	AddBlacklistTokenErrCode = 12006
)

func NewBizErr(err error, errCode int64, errMsg string) error {
	return errors.Wrap(err, fmt.Sprintf("err_code=%d, err_msg=%s", errCode, errMsg))
}

var (
	TokenVoidErr         = func(err error) error { return NewBizErr(err, TokenVoidErrCode, "token无效") }
	TokenExpiredErr      = func(err error) error { return NewBizErr(err, TokenExpiredErrCode, "token已过期") }
	DeleteTokenErr       = func(err error) error { return NewBizErr(err, DeleteTokenErrCode, "token已被撤销") }
	CreateTokenErr       = func(err error) error { return NewBizErr(err, CreateTokenErrCode, "创建token失败") }
	TokenRevokedErr      = func(err error) error { return NewBizErr(err, TokenRevokedErrCode, "token在黑名单中") }
	AddBlacklistTokenErr = func(err error) error {
		return NewBizErr(err, AddBlacklistTokenErrCode, "添加token到黑名单中失败")
	}
)
