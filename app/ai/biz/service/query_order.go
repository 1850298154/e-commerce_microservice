package service

import (
	"2501YTC/app/ai/errno"
	"context"
	"fmt"
	"time"

	einoTool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"

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

	chatModel, err := pkg.CreateARKModel(s.ctx)
	if err != nil {
		err = errno.CreateChatModelErr(err)
		klog.Error(err)
		return
	}
	searchOrdersTool := tool.GetSearchOrdersTool()

	tools := []einoTool.BaseTool{
		searchOrdersTool,
	}

	input := fmt.Sprintf("根据用户id: %d, %s", req.UserId, req.Content)
	// 获取今天的日期
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")
	persona := `你是一个智能助手，集成在一个管理订单和商品的系统中。你的任务是负责根据用户提供的筛选条件从订单信息中筛选出符合要求的订单。
你可以调用的工具是 SearchOrdersTool，该工具可以返回所有订单信息，订单的格式如下：
[
  {
    "order_id": "12345",
    "user_id": 67890,
    "user_currency": "USD",
    "email": "xxx",
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
        "cost": 9.99
      }
    ],
    "orderState": "completed"
  },
  ...
]
用户可能会提供以下筛选条件：
- 创建时间：例如，查找在特定时间段内（比如2025年2月19日或者昨天）创建的订单，为此你需要知道今天的日期为：` + currentDate + `
	。如果用户想要查找的是相对日期的订单，那么你就需要先计算出用户想要查询的具体是哪个时间段，然后再筛选这个时间段内的订单。
	下面举几个例子来说明如何计算相对日期，假设今天是2025年2月26日，那么昨天就是2025年2月25日，前天（或者两天前）就是2025年2月24日，三天前就是2025年2月23日，x天前就是当前日期的天数减x，此外，一周前就是当前日期的天数减7，一个月前就是当前日期的月份减去1，一年前就是当前日期的年份减1。因此如果用户想要查询的是昨天的订单，那就是查找2025年2月25日的订单。
- 商品名称：例如，查找包含特定商品的订单。如果用户提供了多个商品名称（例如商品A、B、C），那么只要订单中包含这三个商品中的任意一个，即为符合要求的订单。
- 订单状态：例如，查找出订单状态为“paid”的订单。
你的任务是：
1. 调用 SearchOrdersTool 获取所有订单信息。
2. 根据用户提供的筛选条件，对订单信息进行筛选。
3. 输出筛选后的订单信息，格式与原始订单信息相同。
请注意：
1. 如果用户提供了多个商品名称，订单中只要包含任意一个商品名称，即为符合要求的订单。
2. 输出的订单信息应与原始订单信息格式一致，并且只输出json形式的数据即可，不要有多余的文字输出。
3. 输出的订单信息一定来自于工具调用中返回的订单信息，不要自己捏造数据。
4. 如果没有符合要求的订单，则输出“[]”。
`

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

	// 直接打印
	orderList, err := pkg.ConvertToAiOrderViewArray(sr.Content)
	if err != nil {
		err = errno.ConvertToAiOrderViewErr(err)
		klog.Error(err)
		return nil, err
	}
	klog.Infof("**********===== finished =====***************")

	return &ai.OrderQueryResp{Order: orderList}, nil
}
