package nodes

import (
	"context"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

var (
	nodeTools = []llms.Tool{}
)

func ToolNode(nodeTools []tools.Tool) graph.NodeFunction {
	return func(ctx context.Context, state []llms.MessageContent, options graph.Options) ([]llms.MessageContent, error) {
		options.CallbackHandler.HandleNodeStart(ctx, "Execute", state)
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
						options.CallbackHandler.HandleNodeEnd(ctx, "Execute", state)
						return state, nil
					}
				}
			}
		}

		options.CallbackHandler.HandleNodeEnd(ctx, "Execute", state)
		return state, nil
	}
}
