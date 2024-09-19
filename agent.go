package main

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

type Agent struct {
	LLM      llms.Model
	Tools    []tools.Tool
	callback *callbacks.AgentFinalStreamHandler
	Executor *agents.Executor
}

func NewAgent(config AppConfig) *Agent {
	agent := &Agent{}

	agent.callback = callbacks.NewFinalStreamHandler()

	agent.LLM = config.AgentLLM()

	mainAgent := agents.NewOneShotAgent(
		agent.LLM,
		agent.Tools,
		agents.WithCallbacksHandler(agent.callback),
	)

	agent.Executor = agents.NewExecutor(
		mainAgent,
		agent.Tools,
	)

	return agent

}

func (agent *Agent) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte)) {
	agent.callback.ReadFromEgress(ctx, callback)
}

func (agent *Agent) Run(ctx context.Context, input string) error {
	log.Println("Agent running with input:", input)
	_, err := chains.Run(ctx, agent.Executor, input)
	if err != nil {
		return err
	}

	return nil
}
