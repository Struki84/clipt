package main

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/tools"
)

func main() {

	input := "What is the meaning of life?"

	config := NewConfig()

	ctx := context.Background()

	callback := callbacks.NewFinalStreamHandler()
	callback.ReadFromEgress(ctx, func(ctx context.Context, chunk []byte) {
		fmt.Print(string(chunk))
	})

	agent := agents.NewOneShotAgent(
		config.AgentLLM(),
		[]tools.Tool{},
		agents.WithCallbacksHandler(callback),
	)

	executor := agents.NewExecutor(
		agent,
		[]tools.Tool{},
	)

	_, err := chains.Run(ctx, executor, input)
	if err != nil {
		fmt.Println(err)
	}
}
