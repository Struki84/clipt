package nodes

import (
	"context"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/tmc/langchaingo/llms"
)

func AgentNode(llm llms.Model, functions []llms.Tool) graph.NodeFunction {
	return func(ctx context.Context, state []llms.MessageContent, options graph.Options) ([]llms.MessageContent, error) {
		options.CallbackHandler.HandleNodeStart(ctx, "Agent", state)

		streamFunc := func(ctx context.Context, chunk []byte) error {
			options.CallbackHandler.HandleNodeStream(ctx, "Agent", chunk)
			return nil
		}

		response, err := llm.GenerateContent(ctx, state, llms.WithTools(functions), llms.WithStreamingFunc(streamFunc))
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
		options.CallbackHandler.HandleNodeEnd(ctx, "Agent", state)
		return state, nil
	}
}
