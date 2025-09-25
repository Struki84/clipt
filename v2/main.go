package v2

import (
	"context"
	"fmt"
	"os"

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
		fmt.Println(err)
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
	return "GPT-4o by OpenAI"
}

func (model TestModel) GetChatHistory() []string {
	return []string{}
}

func (model TestModel) Run(ctx context.Context, input string) error {
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

func (model TestModel) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte) error) {
	model.streamHandler = callback
}

var cliptCmd = &cobra.Command{
	Use:   "chrom",
	Short: "CLI RAG tool for interacting with LLMs.",
	Long:  "",
}

func init() {

	chatCmd := &cobra.Command{
		Use:   "chat",
		Short: "Run chat UI",
		Run: func(cmd *cobra.Command, args []string) {
			model := NewTestModel()
			ShowChatViewLight(model)
		},
	}

	cliptCmd.AddCommand(chatCmd)

}

func main() {
	if err := cliptCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
