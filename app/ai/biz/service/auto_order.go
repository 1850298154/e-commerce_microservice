package service

import (
	"2501YTC/app/ai/infra/rpc"
	"2501YTC/app/ai/pkg"
	autoOrderTool "2501YTC/app/ai/pkg/tool"
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
	"fmt"
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

	//chatModel := pkg.CreateDeepSeekModel(s.ctx)
	chatModel := pkg.CreateARKModel(s.ctx)
	searchProductTool := autoOrderTool.GetSearchProductTool()
	addToCartTool := autoOrderTool.GetAddToCartTool()
	checkoutTool := autoOrderTool.GetCheckoutTool()
	tools := []einoTool.BaseTool{
		searchProductTool,
		addToCartTool,
		checkoutTool,
	}

	persona := `你是一个帮助用户搜索商品，并且下单的助手，根据用户的需要，查询商品信息，并将查到的商品加入到购物车，等商品都加入到购物车后，进行结算。
为了完成这项任务，你需要调用提供的工具，按照以下步骤一个个进行：
1. 根据商品名称查询商品
2. 根据搜寻到的商品ID（就是搜出来的商品里面的id字段）将商品添加到购物车
3. 针对用户想要购买的每种商品，分别进行前两个步骤
4. 将用户想要购买的所有商品添加好购物车后，对用户购物车的商品进行结算，并返回创建好的订单信息
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

	fmt.Println(sr.Content)
	klog.Infof("\n\n===== finished =====\n")

	return nil, nil
}
