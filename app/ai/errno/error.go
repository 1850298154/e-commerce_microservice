package errno

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// System Code
	ServiceErrCode              = 90001
	ParamErrCode                = 90002
	CreateAgentErrCode          = 90003
	StreamErrCode               = 90004
	ConvertToAiOrderViewErrCode = 90005
	CreateChatModelErrCode      = 90006
)

func NewBizErr(err error, errCode int64, errMsg string) error {
	return errors.Wrap(err, fmt.Sprintf("err_code=%d, err_msg=%s", errCode, errMsg))
}

var (
	ServiceErr              = func(err error) error { return NewBizErr(err, ServiceErrCode, "服务没有重新启动") }
	ParamErr                = func(err error) error { return NewBizErr(err, ParamErrCode, "无效的参数请求") }
	CreateAgentErr          = func(err error) error { return NewBizErr(err, CreateAgentErrCode, "创建智能体失败") }
	StreamErr               = func(err error) error { return NewBizErr(err, StreamErrCode, "流式转换失败") }
	ConvertToAiOrderViewErr = func(err error) error { return NewBizErr(err, ConvertToAiOrderViewErrCode, "视图转换失败") }
	CreateChatModelErr      = func(err error) error { return NewBizErr(err, CreateChatModelErrCode, "创建chat模型失败") }
)
