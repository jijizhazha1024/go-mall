package gpt

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type Gpt struct {
	modelID string
	apiKey  string
}

func NewGpt(apiKey, modelID string) *Gpt {
	return &Gpt{
		apiKey:  apiKey,
		modelID: modelID,
	}
}

func (g *Gpt) ChatWithModel(ctx context.Context,
	systemMsg, userMsg string) (string, error) {
	client := arkruntime.NewClientWithApiKey(g.apiKey)
	req := model.ChatCompletionRequest{
		Model: g.modelID,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(systemMsg),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(userMsg),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("standard chat error: %w", err)
	}
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != nil {
		return *resp.Choices[0].Message.Content.StringValue, nil
	}
	return "", fmt.Errorf("no response from the model")
}
