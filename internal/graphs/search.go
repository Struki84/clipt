package graphs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/struki84/clipt/internal/tools/google"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/duckduckgo"
)

var (
	searchPrimer = `You are an agent that has access to a DuckDuckGo and Google search engine.
	Please provide the user with the information they are looking for by using the search tools provided.`

	searchTools = []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "secondarySearch",
				Description: "Performs DuckDuckGo web search, use this search tool only if primary search fails.",
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
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "primarySearch",
				Description: "Performs google web search via serpapi. Use this search tool as primary tool.",
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

var _ tools.Tool = &WebSearchTool{}

type WebSearchTool struct {
	workflow *graph.Runnable
}

func NewWebSearchTool(llm llms.Model) *WebSearchTool {
	return &WebSearchTool{
		workflow: SearchGraph(llm),
	}
}

func (search *WebSearchTool) Name() string {
	return "WebSearch"
}

func (search *WebSearchTool) Description() string {
	return "Performs web search using Google and DuckDuckGo, will resolve to DuckDuckGo if Google is unavailable."
}

func (search *WebSearchTool) Call(ctx context.Context, input string) (string, error) {
	fmt.Println("Performing web search...")
	log.Printf("Input: %s", input)

	initialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, searchPrimer),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	response, err := search.workflow.Invoke(ctx, initialState)
	if err != nil {
		return "", err
	}

	lastMsg := response[len(response)-1].Parts[0].(llms.TextContent).Text

	return lastMsg, nil
}

func SearchGraph(llm llms.Model) *graph.Runnable {

	agent := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		response, err := llm.GenerateContent(ctx, state, llms.WithTools(searchTools))
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

	search := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			if toolCall, ok := part.(llms.ToolCall); ok {
				var args struct {
					Query string `json:"query"`
				}

				if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
					return state, err
				}

				var toolResponse string
				if toolCall.FunctionCall.Name == "primarySearch" {
					apiKey := os.Getenv("SERPAPI_API_KEY")

					google, err := google.New(apiKey, 5)
					if err != nil {
						log.Printf("search error: %v", err)
						return state, err
					}

					toolResponse, err = google.Call(ctx, args.Query)
					if err != nil {
						log.Printf("search error: %v", err)
						return state, err
					}
				}

				if toolCall.FunctionCall.Name == "secondarySearch" {
					search, err := duckduckgo.New(5, duckduckgo.DefaultUserAgent)
					if err != nil {
						log.Printf("search error: %v", err)
						return state, err
					}

					toolResponse, err = search.Call(ctx, args.Query)
					if err != nil {
						log.Printf("search error: %v", err)
						return state, err
					}
				}

				msg := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolCall.ID,
							Name:       toolCall.FunctionCall.Name,
							Content:    toolResponse,
						},
					},
				}

				state = append(state, msg)
			}
		}

		return state, nil
	}

	shouldSearch := func(ctx context.Context, state []llms.MessageContent) string {
		lastMsg := state[len(state)-1]
		for _, part := range lastMsg.Parts {
			if _, ok := part.(llms.ToolCall); ok {
				return "search"
			}
		}

		return graph.END
	}

	workflow := graph.NewMessageGraph()

	workflow.AddNode("agent", agent)
	workflow.AddNode("search", search)

	workflow.SetEntryPoint("agent")
	workflow.AddConditionalEdge("agent", shouldSearch)
	workflow.AddEdge("search", "agent")

	app, err := workflow.Compile()
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}

	return app
}
