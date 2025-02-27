package service

import (
	"2501YTC/app/ai/errno"
	"context"
	"fmt"

	"2501YTC/app/ai/infra/rpc"
	"2501YTC/app/ai/pkg"
	autoOrderTool "2501YTC/app/ai/pkg/tool"
	ai "2501YTC/rpc_gen/kitex_gen/ai"

	einoTool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"
)

type AutoOrderService struct {
	ctx context.Context
} // NewAutoOrderService new AutoOrderService
func NewAutoOrderService(ctx context.Context) *AutoOrderService {
	return &AutoOrderService{ctx: ctx}
}

// Run create note info
func (s *AutoOrderService) Run(req *ai.AutoOrderReq) (resp *ai.AutoOrderResp, err error) {
	// Finish your business logic.
	rpc.InitClient()

	// chatModel := pkg.CreateDeepSeekModel(s.ctx)
	chatModel, err := pkg.CreateARKModel(s.ctx)
	if err != nil {
		err = errno.CreateChatModelErr(err)
		klog.Error(err)
		return
	}
	searchProductTool := autoOrderTool.GetSearchProductTool()
	addToCartTool := autoOrderTool.GetAddToCartTool()
	checkoutTool := autoOrderTool.GetCheckoutTool()
	tools := []einoTool.BaseTool{
		searchProductTool,
		addToCartTool,
		checkoutTool,
	}

	persona := `你是一个帮助用户搜索商品，并且下单的助手，根据用户的需要，查询商品信息，并将查到的商品加入到购物车，等商品都加入到购物车后，进行结算。注意按照用户输入的商品数量进行下单！
请将下单后的订单信息按照json对象的形式进行返回，例如：
		[{
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
				"cost": 19.99
				},
				{
				"product_id": 2,
				"product_name": "Product B",
				"quantity": 1,
				"cost": "9.99"
				}
			],
			"orderState": "placed"
		}]
注意，只返回json形式的数据即可，不要有多余的文字输出，如果没有创建订单，就输出“{}”！
`

	input := fmt.Sprintf("为user_id是%d的用户%s", req.UserId, req.Content)
	ragent, err := react.NewAgent(s.ctx, &react.AgentConfig{
		Model: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},

		MessageModifier: react.NewPersonaModifier(persona),
	})
	if err != nil {
		err = errno.CreateAgentErr(err)
		klog.Error(err)
		return
	}

	sr, err := ragent.Generate(s.ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: input,
		},
	}, agent.WithComposeOptions(compose.WithCallbacks(&pkg.LoggerCallback{})))
	if err != nil {
		err = errno.StreamErr(err)
		klog.Error(err)
		return
	}
	klog.Infof("===== start streaming =====\n\n")
	order, err := pkg.ConvertToAiOrderView(sr.Content)
	if err != nil {
		err = errno.ConvertToAiOrderViewErr(err)
		klog.Error(err)
		return nil, err
	}
	klog.Infof("===== finished =====\n")

	return &ai.AutoOrderResp{Order: order}, nil
}
