package nodes

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

var (
	nodeTools = []llms.Tool{}
)

func ToolNode(nodeTools []tools.Tool) func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
	return func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		fmt.Println("=================== ToolNode ===================")

		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			if toolCall, ok := part.(llms.ToolCall); ok {

				for _, tool := range nodeTools {

					if tool.Name() == toolCall.FunctionCall.Name {
						toolResonse, err := tool.Call(ctx, toolCall.FunctionCall.Arguments)
						if err != nil {
							return state, err
						}

						log.Println("Tool response:", toolResonse)

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
					}
				}
			}
		}
		return state, nil
	}
}
