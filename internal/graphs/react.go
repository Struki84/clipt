package graphs

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools/duckduckgo"
)

var (
	primer = `
	You are a ReAct agent with access to a DuckDuckGo search tool.
	Reason step-by-step to answer the user's query.
	Use the 'search' tool when needed.
	When you have the final answer, end your response with '[FINISH]' on a new line.`

	tools = []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "search",
				Description: "Performs DuckDuckGo web search",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query",
						},
					},
				},
			},
		},
	}
)

func ReactGraph(ctx context.Context, input string) {
	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		log.Fatalf("failed to create LLM: %v", err)
		return
	}

	reason := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		prompt := llms.TextParts(llms.ChatMessageTypeSystem, "Reason step-by-step about the next action to achieve the user's goal based on the current state.")
		state = append(state, prompt)

		resp, err := llm.GenerateContent(ctx, state)
		if err != nil {
			return state, err
		}

		state = append(state, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

		return state, nil
	}

	decision := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		resp, err := llm.GenerateContent(ctx, state, llms.WithTools(tools))
		if err != nil {
			return state, err
		}

		msg := llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content)

		if len(resp.Choices[0].ToolCalls) > 0 {
			for _, toolCall := range resp.Choices[0].ToolCalls {
				msg.Parts = append(msg.Parts, toolCall)
			}
		}

		state = append(state, msg)

		return state, nil
	}

	execute := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			toolCall, ok := part.(llms.ToolCall)

			if ok && toolCall.FunctionCall.Name == "search" {
				var args struct {
					Query string `json:"query"`
				}

				if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
					return state, err
				}

				searchTool, err := duckduckgo.New(1, duckduckgo.DefaultUserAgent)
				if err != nil {
					return state, err
				}

				result, err := searchTool.Call(ctx, args.Query)
				if err != nil {
					return state, err
				}

				msg := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolCall.ID,
							Name:       toolCall.FunctionCall.Name,
							Content:    result,
						},
					},
				}

				return append(state, msg), nil
			}
		}
		return state, nil
	}

	observe := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		prompt := llms.TextParts(llms.ChatMessageTypeSystem, "Review the current state and decide if the user's goal is met. If so, end your response with '[FINISH]' on a new line. If not, suggest the next step.")
		state = append(state, prompt)

		resp, err := llm.GenerateContent(ctx, state)
		if err != nil {
			return state, err
		}

		return append(state, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content)), nil
	}

	shouldAct := func(ctx context.Context, state []llms.MessageContent) string {
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			if toolCall, ok := part.(llms.ToolCall); ok && toolCall.FunctionCall.Name == "search" {
				return "execute"
			}
		}

		return "observe"
	}

	shouldContinue := func(ctx context.Context, state []llms.MessageContent) string {
		lastMsg := state[len(state)-1]

		textContent, ok := lastMsg.Parts[0].(llms.TextContent)

		if ok && strings.Contains(textContent.Text, "[FINISH]") {
			return graph.END
		}

		return "reason"
	}

	workflow := graph.NewMessageGraph()

	workflow.AddNode("reason", reason)
	workflow.AddNode("decide", decision)
	workflow.AddNode("execute", execute)
	workflow.AddNode("observe", observe)

	workflow.SetEntryPoint("reason")
	workflow.AddEdge("reason", "decide")
	workflow.AddConditionalEdge("decide", shouldAct)
	workflow.AddEdge("execute", "observe")
	workflow.AddConditionalEdge("observe", shouldContinue)

	initialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, primer),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	app, err := workflow.Compile()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
		return
	}

	response, err := app.Invoke(ctx, initialState)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	for i, msg := range response {
		log.Printf("Step %d: %v", i, msg.Parts)
	}

	lastMsg := response[len(response)-1]
	log.Printf("last msg: %v", lastMsg.Parts[0])
}
