package graphs

import (
	"context"
	"log"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/struki84/clipt/internal/graphs/nodes"
	"github.com/struki84/clipt/internal/tools/library"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

var (
	libraryPrimer = `You are a agent that can read Office and PDF documents.`

	libraryTools = []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "ReadPDFFile",
				Description: "Reads a PDF file.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query",
						},
						"file": map[string]any{
							"type":        "string",
							"description": "The full filename",
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "ReadOfficeFile",
				Description: "Reads an Office file.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The search query",
						},
						"file": map[string]any{
							"type":        "string",
							"description": "The full filename",
						},
					},
				},
			},
		},
	}
)

var _ tools.Tool = &LibraryTool{}

type LibraryTool struct {
	workflow *graph.Runnable
}

func NewLibraryTool(llm llms.Model) *LibraryTool {
	return &LibraryTool{
		workflow: LibraryGraph(llm),
	}
}

func (tool *LibraryTool) Name() string {
	return "ReadFile"
}

func (tool *LibraryTool) Description() string {
	return "Reads a file."
}

func (tool *LibraryTool) Call(ctx context.Context, input string) (string, error) {
	log.Println("Reading a file...")
	log.Printf("Input: %s", input)

	initialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, libraryPrimer),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}

	response, err := tool.workflow.Invoke(ctx, initialState)
	if err != nil {
		return "", err
	}

	lastMsg := response[len(response)-1].Parts[0].(llms.TextContent).Text

	return lastMsg, nil
}

func LibraryGraph(llm llms.Model) *graph.Runnable {

	read := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			if toolCall, ok := part.(llms.ToolCall); ok {

				var toolResponse string
				if toolCall.FunctionCall.Name == "ReadPDFFile" {
					tool, err := library.NewPDFReaderTool(library.WithModel(llm))
					if err != nil {
						return state, err
					}

					toolResponse, err = tool.Call(ctx, toolCall.FunctionCall.Arguments)
					if err != nil {
						return state, err
					}
				}

				if toolCall.FunctionCall.Name == "ReadOfficeFile" {
					tool, err := library.NewOfficeTool(library.WithModel(llm))
					if err != nil {
						return state, err
					}

					toolResponse, err = tool.Call(ctx, toolCall.FunctionCall.Arguments)
					if err != nil {
						return state, err
					}
				}

				msg := llms.MessageContent{
					Role: llms.ChatMessageTypeAI,
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

	shouldRead := func(ctx context.Context, state []llms.MessageContent) string {
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			if _, ok := part.(llms.ToolCall); ok {
				return "read"
			}
		}

		return graph.END
	}

	workflow := graph.NewMessageGraph()
	workflow.AddNode("agent", nodes.AgentNode(llm, libraryTools))
	workflow.AddNode("read", read)

	workflow.SetEntryPoint("agent")
	workflow.AddConditionalEdge("agent", shouldRead)
	workflow.AddEdge("read", "agent")

	app, err := workflow.Compile()
	if err != nil {
		log.Println("Error compiling graph:", err)
		return nil
	}

	return app
}
