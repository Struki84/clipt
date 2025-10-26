package providers

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAI struct {
	LLM           *openai.LLM
	streamHandler func(ctx context.Context, chunk []byte) error
	currentModel  string
}

func NewOpenAI(model string) *OpenAI {
	llm, err := openai.New(openai.WithModel(model))
	if err != nil {
		fmt.Println("Can't create model:", err)
		return nil
	}
	return &OpenAI{
		LLM:          llm,
		currentModel: model,
	}
}

func (model *OpenAI) Type() string {
	return "LLM"
}

func (model *OpenAI) Name() string {
	return model.currentModel
}

func (model *OpenAI) Description() string {
	return "GPT-4o by OpenAI"
}

func (model *OpenAI) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error) {
	model.streamHandler = callback
}

func (model *OpenAI) Run(ctx context.Context, input string) error {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant!"),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	_, err := model.LLM.GenerateContent(ctx, content, llms.WithStreamingFunc(model.streamHandler))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
