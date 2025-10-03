package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type TestModel struct {
	LLM           *openai.LLM
	streamHandler func(ctx context.Context, chunk []byte) error
}

func NewTestModel() *TestModel {
	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		fmt.Println("Can't create model:", err)
		return nil
	}
	return &TestModel{
		LLM: llm,
	}
}

func (model *TestModel) Type() string {
	return "LLM"
}

func (model *TestModel) Name() string {
	return "GPT-4o"
}

func (model *TestModel) Description() string {
	return "GPT-4o by OpenAI"
}

func (model *TestModel) GetChatHistory(sessionID string) ChatHistory {
	return nil
}

func (model *TestModel) Run(ctx context.Context, input string) error {
	log.Printf("running with input: %v", input)
	log.Printf("stream handler is nil: %v", model.streamHandler == nil)

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant!"),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	completion, err := model.LLM.GenerateContent(ctx, content, llms.WithStreamingFunc(model.streamHandler))
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Printf("Completion: %v", completion.Choices[0].Content)

	return nil
}

func (model *TestModel) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error) {
	log.Printf("Set stream")
	model.streamHandler = callback
}
