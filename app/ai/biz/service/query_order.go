package service

import (
	"context"
	"encoding/json"
	"fmt"
	einoTool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"time"

	"2501YTC/app/ai/infra/rpc"
	"2501YTC/app/ai/pkg"
	"2501YTC/app/ai/pkg/tool"
	ai "2501YTC/rpc_gen/kitex_gen/ai"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"
)

type QueryOrderService struct {
	ctx context.Context
} // NewQueryOrderService new QueryOrderService
func NewQueryOrderService(ctx context.Context) *QueryOrderService {
	return &QueryOrderService{ctx: ctx}
}

// Run create note info
func (s *QueryOrderService) Run(req *ai.OrderQueryReq) (resp *ai.OrderQueryResp, err error) {
	// Finish your business logic.
	rpc.InitClient()

	//chatModel := pkg.CreateDeepSeekModel(s.ctx)
	chatModel := pkg.CreateARKModel(s.ctx)
	//searchOrderByProductNameTool := tool.GetSearchOrderByProductNameTool()
	searchOrdersTool := tool.GetSearchOrdersTool()

	tools := []einoTool.BaseTool{
		//searchOrderByProductNameTool,
		searchOrdersTool,
	}

	input := fmt.Sprintf("根据用户id: %d, %s", req.UserId, req.Content)
	// 获取今天的日期
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")
	persona := "你是一个智能助手，集成在一个管理订单和商品的系统中。你的任务是帮助用户使用提供的函数调用工具根查询用户的所有订单并根据特定条件从中筛选订单信息，如果用户输入的是商品名称，则需要筛选出所有包含用户指定商品的订单。为了完成这个任务，你需要知道今天的日期为" + currentDate + "。" +
		"请将符合条件的订单信息按照json对象的形式进行返回，例如：\n" +
		`[{
		  "order_id": "12345",
		  "user_id": 67890,
		  "user_currency": "USD",
		  "email": "user@example.com",
		  "created_at": "2023-10-01T12:34:56Z",
		  "order_items": [
			{
			  "product_id": 1,
			  "product_name": "Product A",
			  "quantity": 2,
			  "cost": "19.99"
			},
			{
			  "product_id": 2,
			  "product_name": "Product B",
			  "quantity": 1,
			  "cost": "9.99"
			}
		  ],
		  "orderState": "completed"
		},...]
注意，只返回json形式的数据即可，不要有多余的文字输出，如果未找到符合条件的订单信息，就输出“[]”！`

	ragent, err := react.NewAgent(s.ctx, &react.AgentConfig{
		Model: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},

		MessageModifier: react.NewPersonaModifier(persona),
	})
	if err != nil {
		klog.Errorf("failed to create agent: %v", err)
		return
	}

	sr, err := ragent.Generate(s.ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: input,
		},
	}, agent.WithComposeOptions(compose.WithCallbacks(&pkg.LoggerCallback{})))
	if err != nil {
		klog.Errorf("failed to stream: %v", err)
		return
	}

	klog.Infof("\n\n===== start streaming =====\n\n")

	// 直接打印
	respContent := sr.Content
	var orderList []*ai.OrderResult
	err = json.Unmarshal([]byte(respContent), &orderList)
	fmt.Println(orderList)
	klog.Infof("**********===== finished =====***************")

	return &ai.OrderQueryResp{Order: orderList}, nil
}
