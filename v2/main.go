package v2

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type TestModel struct {
	LLM           *openai.LLM
	streamHandler func(ctx context.Context, chunk []byte) error
}

func NewTestModel() TestModel {
	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		log.Fatalf("Failed to load the model")
		return TestModel{}
	}

	return TestModel{
		LLM: llm,
	}
}

func (model TestModel) Type() string {
	return "LLM"
}

func (model TestModel) Name() string {
	return "GPT-4o"
}

func (model TestModel) Description() string {
	return "GTP-4o by OpenAI"
}

func (model TestModel) Run(ctx context.Context, input string) error {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are a helpful assistant!"),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	_, err := model.LLM.GenerateContent(ctx, content, llms.WithStreamingFunc(model.streamHandler))
	if err != nil {
		log.Fatalf("Failed to generate content!")
		return err
	}
	return nil
}

func (model TestModel) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error) {
	model.streamHandler = callback
}

var chromCmd = &cobra.Command{
	Use:   "chrom",
	Short: "CLI RAG tool for interacting with LLMs.",
	Long:  "",
}

func init() {

}

func main() {

}
