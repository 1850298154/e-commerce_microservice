package templates

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"time"
)

func CreateTemplate() prompt.ChatTemplate {
	// 获取当前的日期和时间
	currentTime := time.Now()
	// 格式化并只获取日期部分
	currentDate := currentTime.Format("2006-01-02")

	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你的目标是帮助商城的用户完成想要的指令。为了完成指令，你需要知道今天的日期为"+currentDate),

		// 用户消息模板
		schema.UserMessage("指令: {task}"),
	)
}
