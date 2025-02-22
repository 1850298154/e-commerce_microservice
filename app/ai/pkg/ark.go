package pkg

import (
	"context"
	"github.com/cloudwego/eino/components/model"
	"github.com/joho/godotenv"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/kitex/pkg/klog"
)

func CreateARKModel(ctx context.Context) model.ChatModel {
	_ = godotenv.Load("../../.env")
	arkApi := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("ARK_MODEL_NAME")
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: arkApi,
		Model:  modelName,
	})

	if err != nil {
		klog.Errorf("failed to create chat model: %v", err)
		return nil
	}
	return chatModel
}
