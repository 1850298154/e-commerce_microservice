package templates

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func CreateTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你的目标是帮助商城的用户完成想要的指令。"),

		// 用户消息模板
		schema.UserMessage("指令: {task}"),
	)
}
