package providers

import (
	"context"
	"fmt"
	"log"

	"github.com/struki84/clipt/storage"
	"github.com/struki84/clipt/tui"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAI struct {
	LLM           *openai.LLM
	streamHandler func(ctx context.Context, chunk []byte) error
	currentModel  string
	storage       storage.SQLite
}

func NewOpenAI(model string, storage storage.SQLite) *OpenAI {
	llm, err := openai.New(openai.WithModel(model))
	if err != nil {
		fmt.Println("Can't create model:", err)
		return nil
	}
	return &OpenAI{
		LLM:          llm,
		currentModel: model,
		storage:      storage,
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

func (model *OpenAI) Run(ctx context.Context, input string, session tui.ChatSession) error {
	buffer, err := model.storage.LoadSessionMsgs(session.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant!"),
		llms.TextParts(llms.ChatMessageTypeSystem, "CHAT HISTORY: \n"+buffer),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	response, err := model.LLM.GenerateContent(ctx, content, llms.WithStreamingFunc(model.streamHandler))
	if err != nil {
		fmt.Println(err)
		return err
	}

	text := response.Choices[0].Content
	err = model.storage.SaveSessionMsg(session.ID, input, text)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
