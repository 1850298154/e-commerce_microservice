package utils

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
)

func MustHandleError(err error) {
	if err != nil {
		hlog.Fatal(err)
	}
}

const ()

func NewBizErr(errCode int64, errMsg string) error {
	return errors.New(fmt.Sprintf("err_code=%d, err_msg=%s", errCode, errMsg))
}
