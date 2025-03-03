package utils

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func MustHandleError(err error) {
	if err != nil {
		hlog.Fatal(err)
	}
}

func NewBizErr(errCode int64, errMsg string) error {
	return fmt.Errorf("err_code=%d, err_msg=%s", errCode, errMsg)
}
