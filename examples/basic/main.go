package main

import (
	"context"
	"fmt"

	"github.com/struki84/clipt"
	"github.com/struki84/clipt/tui"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/openai"
)

type AgentProvider struct {
}

func (agent AgentProvider) Type() string {
	return "Agent"
}

func (agent AgentProvider) Name() string {
	return "ReAct"
}

func (agent AgentProvider) Description() string {
	return "Basic React agent"
}

func (agent AgentProvider) ChatHistory(sessionID string) tui.ChatHistory {
	return nil
}

func (agent AgentProvider) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error) {

}

func (agent AgentProvider) Run(ctx context.Context, input string) error {
	return nil
}

type TestProvider2 struct {
	LLM           *anthropic.LLM
	streamHandler func(ctx context.Context, chunk []byte) error
}

func NewTestProvider2() *TestProvider2 {
	llm, err := anthropic.New()
	if err != nil {
		fmt.Println("Can't create model:", err)
		return nil
	}
	return &TestProvider2{
		LLM: llm,
	}
}

func (model *TestProvider2) Type() string {
	return "LLM"
}

func (model *TestProvider2) Name() string {
	return "Claude"
}

func (model *TestProvider2) Description() string {
	return "Claude by Anthropic"
}

func (model *TestProvider2) ChatHistory(sessionID string) tui.ChatHistory {
	return nil
}

func (model *TestProvider2) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error) {
	model.streamHandler = callback
}

func (model *TestProvider2) Run(ctx context.Context, input string) error {
	return nil
}

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

	clipt.SetProviders([]tui.ChatProvider{provider, NewTestProvider2(), AgentProvider{}})
	clipt.Render(provider)
}
