package apiErr

import (
	"errors"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

const (
	// System Code
	ServiceErrCode = iota + 10000

	// Business Code
	ProductIDRequired = iota + 40000
)

var (
	ServiceErr = kerrors.NewBizStatusError(ServiceErrCode, "服务异常")

	ProductIDRequiredErr = kerrors.NewBizStatusError(ProductIDRequired, "商品ID不能为空")
)

// ConvertErr convert error to	kerrors.BizStatusErrorIface
func ConvertErr(err error) kerrors.BizStatusErrorIface {
	var e kerrors.BizStatusErrorIface
	if errors.As(err, &e) {
		return e
	}

	s := kerrors.NewBizStatusErrorWithExtra(ServiceErrCode, "服务异常", map[string]string{"error": err.Error()})
	return s
}
