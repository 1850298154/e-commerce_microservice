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
	ProductNameRequired
	ProductPriceInvalid
	PageRequired
	PageSizeRequired
	ImageDataRequired
	FileNameRequired
)

var (
	ServiceErr = kerrors.NewBizStatusError(ServiceErrCode, "服务异常")

	ProductIDRequiredErr   = kerrors.NewBizStatusError(ProductIDRequired, "商品ID不能为空")
	ProductNameRequiredErr = kerrors.NewBizStatusError(ProductNameRequired, "商品名称不能为空")
	ProductPriceInvalidErr = kerrors.NewBizStatusError(ProductPriceInvalid, "商品价格不能为负")
	PageRequiredErr        = kerrors.NewBizStatusError(PageRequired, "页码不能为空")
	PageSizeRequiredErr    = kerrors.NewBizStatusError(PageSizeRequired, "每页数量不能为空")
	ImageDataRequiredErr   = kerrors.NewBizStatusError(ImageDataRequired, "图片数据不能为空")
	FileNameRequiredErr    = kerrors.NewBizStatusError(FileNameRequired, "文件名不能为空")
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
