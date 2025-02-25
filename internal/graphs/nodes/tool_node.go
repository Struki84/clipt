package nodes

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

var (
	nodeTools = []llms.Tool{}
)

func ToolNode(nodeTools []tools.Tool) func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
	return func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		fmt.Println("=================== Tool Node ===================")
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			toolCall, ok := part.(llms.ToolCall)
			if ok {
				for _, tool := range nodeTools {
					if tool.Name() == toolCall.FunctionCall.Name {
						toolResonse, err := tool.Call(ctx, toolCall.FunctionCall.Arguments)
						if err != nil {
							return state, err
						}

						msg := llms.MessageContent{
							Role: llms.ChatMessageTypeTool,
							Parts: []llms.ContentPart{
								llms.ToolCallResponse{
									ToolCallID: toolCall.ID,
									Name:       toolCall.FunctionCall.Name,
									Content:    toolResonse,
								},
							},
						}

						state = append(state, msg)
						return state, nil
					}
				}
			}
		}
		return state, nil
	}
}
