package graphs

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/struki84/clipt/internal/tools/google"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
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
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "primarySearch",
				Description: "Performs google web search via serpapi",
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

func SearchGraph(input string) {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	intialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, ""),
	}

	agent := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		response, err := model.GenerateContent(ctx, state, llms.WithTools(tools))
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
					google, err := google.New("", 5)
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
		return
	}

	intialState = append(
		intialState,
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	)

	response, err := app.Invoke(context.Background(), intialState)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	lastMsg := response[len(response)-1]
	log.Printf("last msg: %v", lastMsg.Parts[0])
}
