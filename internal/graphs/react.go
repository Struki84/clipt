package graphs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/struki84/clipt/internal/graphs/nodes"
	"github.com/struki84/clipt/internal/tools/library"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

var (
	primer = `
	You are a ReAct agent who can serch the web and read files.
	Reason step-by-step to answer the user's query.
	Use the 'WebSearch' tool ONLY when needed.
	Use the 'FileList' tool to list files you can read and search through.
	Use the 'ReadFile tool to read Office and PDF files.'
	`

	reasonPrimer = `
	Reason step-by-step about the next action to achieve the user's goal based on the current state.
	`

	observePrimer = `
	Review the current state and decide if the user's goal is met. 
	If the user's goal is met, construct your final response to the user and end it with '[FINISH]' on a new line. 
	If not, suggest the next step.
	`

	functions = []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "WebSearch",
				Description: "Performs Google web search, will resolve to DuckDuckGo if Google is unavailable.",
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
				Name:        "ListFiles",
				Description: "Lists all files you can read and search.",
				Parameters:  map[string]any{},
			},
		},
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "ReadFile",
				Description: "Use this tool to read and search the contents of Office and PDF files.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "Search query or inquery about the file",
						},
						"file": map[string]any{
							"type":        "string",
							"description": "The file to read",
						},
					},
				},
			},
		},
	}

	graphTools = []tools.Tool{}
)

func ReactGraph(ctx context.Context, input string) {
	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		log.Fatalf("failed to create LLM: %v", err)
		return
	}

	graphTools = append(graphTools,
		NewWebSearchTool(llm),
		NewLibraryTool(llm),
		library.NewFileListTool(),
	)

	reason := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		fmt.Println("=================== Reason ===================")

		prompt := llms.TextParts(llms.ChatMessageTypeSystem, reasonPrimer)
		state = append(state, prompt)

		resp, err := llm.GenerateContent(ctx, state)
		if err != nil {
			return state, err
		}

		state = append(state, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

		log.Println(resp.Choices[0].Content)

		return state, nil
	}

	decide := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		fmt.Println("=================== Act ===================")
		resp, err := llm.GenerateContent(ctx, state, llms.WithTools(functions))
		if err != nil {
			return state, err
		}

		if len(resp.Choices[0].ToolCalls) > 0 {
			log.Println("tool call:", resp.Choices[0].ToolCalls[0])
			msg := llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content)

			for _, toolCall := range resp.Choices[0].ToolCalls {
				msg.Parts = append(msg.Parts, toolCall)
			}

			state = append(state, msg)
		}

		return state, nil
	}

	observe := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		fmt.Println("=================== Observe ===================")

		prompt := llms.TextParts(llms.ChatMessageTypeSystem, observePrimer)
		state = append(state, prompt)

		resp, err := llm.GenerateContent(ctx, state)
		if err != nil {
			return state, err
		}

		log.Println(resp.Choices[0].Content)

		return append(state, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content)), nil
	}

	shouldAct := func(ctx context.Context, state []llms.MessageContent) string {
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			if _, ok := part.(llms.ToolCall); ok {
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
	workflow.AddNode("decide", decide)
	workflow.AddNode("execute", nodes.ToolNode(graphTools))
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

	fmt.Println("=================== Final Response ===================")

	lastMsg := response[len(response)-1]
	fmt.Println(lastMsg.Parts[0].(llms.TextContent).Text)
}
