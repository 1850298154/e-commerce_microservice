package pkg

import (
	"context"
	"os"

	"github.com/cloudwego/eino/components/model"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

func CreateARKModel(ctx context.Context) (model.ChatModel, error) {
	// _ = godotenv.Load("../../.env")
	arkApi := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("ARK_MODEL_NAME")
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: arkApi,
		Model:  modelName,
	})
	if err != nil {
		return nil, err
	}
	return chatModel, nil
}
