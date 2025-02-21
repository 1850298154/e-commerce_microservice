package service

import (
	"2501YTC/app/ai/infra/rpc"
	"2501YTC/app/ai/pkg"
	queryOrderTool "2501YTC/app/ai/pkg/tool"
	ai "2501YTC/rpc_gen/kitex_gen/ai"
	"context"
	"encoding/json"
	"fmt"
	einoTool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"

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
	searchOrderTool := queryOrderTool.GetSearchOrdersTool()
	queryProductTool := queryOrderTool.GetProductTool()

	tools := []einoTool.BaseTool{
		searchOrderTool,
		queryProductTool,
	}

	input := fmt.Sprintf("根据用户id: %d, %s", req.UserId, req.Content)
	persona := `你是一个帮助用户查询订单的助手，根据用户的指示和提供的工具，分析用户的需求，执行相应的工具调用来查询对应的订单信息,并输出json格式的数据。
为了完成此任务，你需要知道以下内容：
1. 今天是2025年2月20日。
2. 如果用户输入了日期信息，这个日期对应订单信息中的创建时间，比如用户输入了2025年2月19日，那么说明用户需要查找创建时间在2025年2月19日的订单
3. 如果用户想要查询购买过某个或某些商品的订单，你需要根据当前给出的所有信息，按照以下步骤一步步进行查询：
	第一步：查询所有订单；
	第二步：从第一步获取的订单信息中提取出product_id；
	第三步：根据product_id查询商品；
	第四步：从商品信息中提取product_name，和用户想要查询的商品名称进行比对，如果是用户要找的商品，则提取出商品信息中的id作为product_id；
	第五步：将第四步提取出来的product_id和第一步中查出来的订单信息中的product_id进行对比，如果某个订单包含第四步的product_id，则将该订单信息进行输出。
`

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

	//defer sr.Close() // remember to close the stream

	klog.Infof("\n\n===== start streaming =====\n\n")

	//for {
	//	msg, err := sr.Recv()
	//	if err != nil {
	//		if errors.Is(err, io.EOF) {
	//			// finish
	//			break
	//		}
	//		// error
	//		klog.Infof("failed to recv: %v", err)
	//		return nil, err
	//	}
	//
	//	// 打字机打印
	//	//fmt.Printf("%v", msg.Content)
	//	//var order []rpcorder.OrderResp
	//
	//}

	// 直接打印
	var orderList []queryOrderTool.SearchOrdersResult
	err = json.Unmarshal([]byte(sr.Content), &orderList)
	fmt.Println(orderList)
	fmt.Println("**********===== finished =====***************")

	return nil, nil
}
