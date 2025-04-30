package graphs

import (
	"context"
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
	Review the current state and first decide if the user's request is met. 
	If the user's goal is met, construct your final response to the user and wrap it with '[FINISH][/FINISH]' tags on a new line. 
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

func ReactGraph(ctx context.Context, input string, callback graph.GraphCallback) {
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

	workflow := graph.NewMessageGraph(graph.WithCallback(callback))

	reason := func(ctx context.Context, state []llms.MessageContent, options graph.Options) ([]llms.MessageContent, error) {
		options.CallbackHandler.HandleNodeStart(ctx, "Reason", state)

		streamFunc := func(ctx context.Context, chunk []byte) error {
			options.CallbackHandler.HandleNodeStream(ctx, "Reason", chunk)
			return nil
		}

		resp, err := llm.GenerateContent(ctx, state, llms.WithStreamingFunc(streamFunc))
		if err != nil {
			return state, err
		}

		state = append(state, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content))

		options.CallbackHandler.HandleNodeEnd(ctx, "Reason", state)
		return state, nil
	}

	act := func(ctx context.Context, state []llms.MessageContent, options graph.Options) ([]llms.MessageContent, error) {
		options.CallbackHandler.HandleNodeStart(ctx, "Act", state)

		stramFunc := func(ctx context.Context, chunk []byte) error {
			options.CallbackHandler.HandleNodeStream(ctx, "Act", chunk)
			return nil
		}

		resp, err := llm.GenerateContent(ctx, state,
			llms.WithTools(functions),
			llms.WithStreamingFunc(stramFunc),
		)

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

		options.CallbackHandler.HandleNodeEnd(ctx, "Act", state)
		return state, nil
	}

	observe := func(ctx context.Context, state []llms.MessageContent, options graph.Options) ([]llms.MessageContent, error) {
		options.CallbackHandler.HandleNodeStart(ctx, "Observe", state)

		prompt := llms.TextParts(llms.ChatMessageTypeSystem, observePrimer)
		state = append(state, prompt)

		streamFunc := func(ctx context.Context, chunk []byte) error {
			options.CallbackHandler.HandleNodeStream(ctx, "Observe", chunk)
			return nil
		}

		resp, err := llm.GenerateContent(ctx, state, llms.WithStreamingFunc(streamFunc))
		if err != nil {
			return state, err
		}

		options.CallbackHandler.HandleNodeEnd(ctx, "Observe", state)
		return append(state, llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content)), nil
	}

	shouldAct := func(ctx context.Context, state []llms.MessageContent, options graph.Options) string {
		callback.HandleEdgeEntry(ctx, "shouldAct", state)
		lastMsg := state[len(state)-1]

		for _, part := range lastMsg.Parts {
			if _, ok := part.(llms.ToolCall); ok {
				callback.HandleEdgeExit(ctx, "shouldAct", state, "execute")
				return "execute"
			}
		}

		callback.HandleEdgeExit(ctx, "shouldAct", state, "observe")
		return "observe"
	}

	shouldContinue := func(ctx context.Context, state []llms.MessageContent, options graph.Options) string {
		callback.HandleEdgeEntry(ctx, "shouldContinue", state)
		lastMsg := state[len(state)-1]

		textContent, ok := lastMsg.Parts[0].(llms.TextContent)

		if ok && strings.Contains(textContent.Text, "[FINISH]") {
			callback.HandleEdgeExit(ctx, "shouldContinue", state, graph.END)
			return graph.END
		}

		callback.HandleEdgeExit(ctx, "shouldContinue", state, "reason")
		return "reason"
	}

	workflow.AddNode("reason", reason)
	workflow.AddNode("act", act)
	workflow.AddNode("execute", nodes.ToolNode(graphTools))
	workflow.AddNode("observe", observe)

	workflow.SetEntryPoint("reason")
	workflow.AddEdge("reason", "act")
	workflow.AddConditionalEdge("act", shouldAct)
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

	_, err = app.Invoke(ctx, initialState)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
}
