package graphs

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/duckduckgo"
	"github.com/tmc/langgraphgo/graph"
	"github.com/tmc/langgraphgo/nodes"
)

func SearchGraph(input string) {

	model, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		panic(err)
	}

	intialState := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "You are an agent that has access to a Duck Duck go search engine. Please provide the user with the information they are looking for by using the search tool provided."),
	}

	ddg, err := duckduckgo.New(3, duckduckgo.DefaultUserAgent)
	if err != nil {
		panic(err)
	}

	agentTools := []llms.Tool{
		{
			Type:     "function",
			Function: ddg.Definition(),
		},
	}

	agent := func(ctx context.Context, state []llms.MessageContent) ([]llms.MessageContent, error) {
		response, err := model.GenerateContent(ctx, state, llms.WithTools(agentTools))
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

	shouldSearch := func(ctx context.Context, state []llms.MessageContent) string {
		lastMsg := state[len(state)-1]
		for _, part := range lastMsg.Parts {
			toolCall, ok := part.(llms.ToolCall)

			if ok && toolCall.FunctionCall.Name == "DuckDuckGoSearch" {
				log.Printf("agent should use search")
				return "search"
			}
		}

		return graph.END
	}

	workflow := graph.NewMessageGraph()

	workflow.AddNode("agent", agent)
	workflow.AddNode("search", nodes.ToolNode([]tools.Tool{ddg}))

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

	lastMsg := response[len(response)-1].Parts[0]
	fmt.Printf("ANSWER: %v", lastMsg)
}
