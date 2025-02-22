package service

import (
	"2501YTC/app/ai/infra/rpc"
	"2501YTC/app/ai/pkg"
	queryOrderTool "2501YTC/app/ai/pkg/tool"
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
	filterTool := queryOrderTool.GetFilterTool()

	tools := []einoTool.BaseTool{
		searchOrderTool,
		queryProductTool,
		filterTool,
	}

	input := fmt.Sprintf("根据用户id: %d, %s", req.UserId, req.Content)
	persona := `你是一个智能助手，集成在一个管理订单和商品的系统中。你的任务是帮助用户使用提供的函数调用工具根据特定条件查询订单。可用的工具有SearchOrdersTool、QueryProductTool和FilterTool。
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
	//var orderList []queryOrderTool.SearchOrdersResult
	//err = json.Unmarshal([]byte(sr.Content), &orderList)
	//fmt.Println(orderList)
	fmt.Println(sr.Content)
	fmt.Println("**********===== finished =====***************")

	return nil, nil
}
