package pkg

import (
	"2501YTC/app/ai/pkg/tool"
	"2501YTC/rpc_gen/kitex_gen/ai"
	"encoding/json"
	"github.com/cloudwego/kitex/pkg/klog"
	"time"
)

func ConvertToAiOrderViewArray(output string) ([]*ai.OrderResult, error) {
	var results []*tool.SearchOrdersResult
	err := json.Unmarshal([]byte(output), &results)
	if err != nil {
		klog.Errorf("failed to unmarshal ai output: %v", err)
		return nil, err
	}
	orderList := make([]*ai.OrderResult, 0)
	var orderItems []*ai.OrderItem
	for _, order := range results {
		orderItems = make([]*ai.OrderItem, 0)
		for _, item := range order.OrderItems {
			orderItems = append(orderItems, &ai.OrderItem{
				ProductId:   item.ProductId,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				Cost:        item.Cost,
			})
		}
		orderList = append(orderList, &ai.OrderResult{
			OrderId:      order.OrderId,
			UserId:       order.UserId,
			UserCurrency: order.UserCurrency,
			Email:        order.Email,
			CreatedAt:    TimeToUnixInt32(order.CreatedAt),
			OrderItems:   orderItems,
			OrderState:   order.OrderState,
		})
	}
	return orderList, nil
}

func ConvertToAiOrderView(output string) (*ai.OrderResult, error) {
	var result *tool.SearchOrdersResult
	err := json.Unmarshal([]byte(output), &result)
	if err != nil {
		klog.Errorf("failed to unmarshal ai output: %v", err)
		return nil, err
	}
	orderItems := make([]*ai.OrderItem, 0)
	for _, item := range result.OrderItems {
		orderItems = append(orderItems, &ai.OrderItem{
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Cost:        item.Cost,
		})
	}
	order := &ai.OrderResult{
		OrderId:      result.OrderId,
		UserId:       result.UserId,
		UserCurrency: result.UserCurrency,
		Email:        result.Email,
		CreatedAt:    TimeToUnixInt32(result.CreatedAt),
		OrderItems:   orderItems,
		OrderState:   result.OrderState,
	}
	return order, nil
}

func TimeToUnixInt32(t time.Time) int32 {
	return int32(t.Unix())
}
