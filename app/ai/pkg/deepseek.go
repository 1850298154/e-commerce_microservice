package pkg

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/kitex/pkg/klog"
)

func CreateDeepSeekModel(ctx context.Context) model.ChatModel {
	// _ = godotenv.Load("../../.env")
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("BASE_URL")

	// 创建并配置 ChatModel
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:  "deepseek-ai/DeepSeek-V2.5",
		APIKey: openAIAPIKey,
		// Temperature: gptr.Of(float32(0.7)),
		BaseURL: baseURL,
	})
	if err != nil {
		klog.Errorf("NewChatModel failed, err=%v", err)
		return nil
	}
	return chatModel
}
