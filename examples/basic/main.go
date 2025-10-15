package main

import (
	"context"
	"fmt"

	"github.com/struki84/clipt/tui"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type TestProvider struct {
	LLM           *openai.LLM
	streamHandler func(ctx context.Context, chunk []byte) error
}

func NewTestProvider() *TestProvider {
	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		fmt.Println("Can't create model:", err)
		return nil
	}
	return &TestProvider{
		LLM: llm,
	}
}

func (model *TestProvider) Type() string {
	return "LLM"
}

func (model *TestProvider) Name() string {
	return "GPT-4o"
}

func (model *TestProvider) Description() string {
	return "GPT-4o by OpenAI"
}

func (model *TestProvider) ChatHistory(sessionID string) tui.ChatHistory {
	return nil
}

func (model *TestProvider) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error) {
	model.streamHandler = callback
}

func (model *TestProvider) Run(ctx context.Context, input string) error {
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

func main() {
	provider := NewTestProvider()

	tui.Render(provider)
}
