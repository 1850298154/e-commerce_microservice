package utils

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetUserIdFromCtx(ctx context.Context) uint32 {
	value := ctx.Value(UserIdKey)
	if value == nil {
		return 0
	}

	userID, ok := value.(uint32)
	if !ok {
		return 0
	}
	return userID
}

func GetUserIdFromReqCtx(c *app.RequestContext) uint32 {
	value, ok := c.Get(UserIdKey)
	if !ok || value == nil {
		hlog.Warnf("GetUserIdFromReqCtxFailed: %v", value)
		return 0
	}

	userID, ok := value.(uint32)
	if !ok {
		hlog.Warnf("GetUserIdFromReqCtxFailed: %v", value)
		return 0
	}
	return userID
}
