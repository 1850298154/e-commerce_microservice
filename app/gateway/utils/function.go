package utils

import "context"

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
