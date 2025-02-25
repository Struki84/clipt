package nodes

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

func AgentNode(llm llms.Model, functions []llms.Tool) func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
	return func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {

		response, err := llm.GenerateContent(ctx, state, llms.WithTools(functions))
		if err != nil {
			return state, err
		}

		msg := llms.TextParts(llms.ChatMessageTypeAI, response.Choices[0].Content)

		if len(response.Choices[0].ToolCalls) > 0 {
			for _, toolCall := range response.Choices[0].ToolCalls {
				msg.Parts = append(msg.Parts, toolCall)
			}
		}

		state = append(state, msg)
		return state, nil
	}
}
