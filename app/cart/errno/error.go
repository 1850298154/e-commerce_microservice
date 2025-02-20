package errno

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// System Code
	CartServiceErrCode = 20001
	CartParamErrCode   = 20002

	// Cart ErrCode
	ProductNotFoundErrCode      = 21001
	InvalidQuantityErrCode      = 21002
	CartEmptyErrCode            = 21003
	ProductAlreadyInCartErrCode = 21004
	InsufficientStockErrCode    = 21005
	AddToCartErrCode            = 21006
	RemoveFromCartErrCode       = 21007
	UpdateCartItemErrCode       = 21008
	GetCartErrCode              = 21009
)

func NewBizErr(err error, errCode int64, errMsg string) error {
	return errors.Wrap(err, fmt.Sprintf("err_code=%d, err_msg=%s", errCode, errMsg))
}

var (
	CartServiceErr          = func(err error) error { return NewBizErr(err, CartServiceErrCode, "购物车服务出现问题") }
	CartParamErr            = func(err error) error { return NewBizErr(err, CartParamErrCode, "购物车参数请求无效") }
	ProductNotFoundErr      = func(err error) error { return NewBizErr(err, ProductNotFoundErrCode, "商品未找到") }
	InvalidQuantityErr      = func(err error) error { return NewBizErr(err, InvalidQuantityErrCode, "商品数量不合法") }
	CartEmptyErr            = func(err error) error { return NewBizErr(err, CartEmptyErrCode, "购物车为空") }
	ProductAlreadyInCartErr = func(err error) error { return NewBizErr(err, ProductAlreadyInCartErrCode, "商品已在购物车中") }
	InsufficientStockErr    = func(err error) error { return NewBizErr(err, InsufficientStockErrCode, "商品库存不足") }
	AddToCartErr            = func(err error) error { return NewBizErr(err, AddToCartErrCode, "添加商品到购物车失败") }
	RemoveFromCartErr       = func(err error) error { return NewBizErr(err, RemoveFromCartErrCode, "从购物车移除商品失败") }
	UpdateCartItemErr       = func(err error) error {
		return NewBizErr(err, UpdateCartItemErrCode, "更新购物车商品信息失败")
	}
	GetCartErr = func(err error) error { return NewBizErr(err, GetCartErrCode, "获取购物车失败") }
)
