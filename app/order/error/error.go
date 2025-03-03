package error

import (
	"fmt"
)

const (
	// User related errors
	ErrInvalidUserId = 18001

	// Order validation errors
	ErrInvalidOrderId    = 18002
	ErrInvalidOrderItems = 18003

	// Order creation errors
	ErrGenerateOrderIdFailed = 18004
	ErrCreateOrderFailed     = 18005
	ErrCreateOrderItemFailed = 18006

	// Order query errors
	ErrGetOrderByUserIdAndOrderIdFailed = 18009
	ErrListOrderByUserIdFailed          = 18012

	// Order modification errors
	ErrUpdateOrderFailed      = 18010
	ErrUpdateOrderItemsFailed = 18011
	ErrCancelOrderFailed      = 18013

	// Queue related errors
	ErrDelayCancelSendToQueueFailed = 18007

	// Database errors
	ErrMySQLTransactionFailed = 18008
)

func NewError(code int, msg string, err error) error {
	if err != nil {
		return fmt.Errorf("code: %d, msg: %s, error: %v", code, msg, err)
	}
	return fmt.Errorf("code: %d, msg: %s", code, msg)
}
