package service

import (
	"context"

	"2501YTC/app/ai/infra/rpc"
	"2501YTC/app/ai/pkg"
	"2501YTC/app/ai/pkg/templates"
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

	// chatTemplate := templates.CreateQueryOrderMessageFromTemplate()
	chatTemplate := templates.CreateTemplate()

	chatModel := pkg.CreateDeepSeekModel(s.ctx)

	tools := tool.NewQueryOrderTools()

	// 获取工具信息, 用于绑定到 ChatModel
	toolInfos, err := tool.GetToolInfo(s.ctx, tools)
	if err != nil {
		return nil, err
	}
	// 将 tool 绑定到 ChatModel
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		klog.Errorf("BindTools failed, err=%v", err)
		return
	}

	// 创建 tool 节点
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		klog.Errorf("NewToolNode failed, err=%v", err)
		return
	}

	// 构建完整的处理链
	chain := compose.NewChain[map[string]any, []*schema.Message]()
	chain.
		AppendChatTemplate(chatTemplate, compose.WithNodeName("chat_template")).
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tool"))

	// 编译并运行 chain
	agent, err := chain.Compile(s.ctx)
	if err != nil {
		klog.Errorf("chain.Compile failed, err=%v", err)
		return
	}

	// 与模型对话
	respMsg, err := agent.Invoke(s.ctx, map[string]any{
		"role": "自动查询订单的助手",
		"task": req.Content,
	},
	)
	if err != nil {
		klog.Errorf("agent.Invoke failed, err=%v", err)
		return
	}

	// respMsg, err := util.Chat(s.ctx, tool, req.Content)
	// 输出结果
	for idx, msg := range respMsg {
		klog.Infof("\n")
		klog.Infof("message %d: %s: %s", idx, msg.Role, msg.Content)
	}
	return
}
