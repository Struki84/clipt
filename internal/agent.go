package internal

import (
	"context"
	"log"

	"github.com/struki84/clipt/config"
	"github.com/struki84/clipt/internal/callbacks"
	mem "github.com/struki84/clipt/internal/memory"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/tools"
)

type Agent struct {
	config   config.AppConfig
	LLM      llms.Model
	Tools    []tools.Tool
	callback *callbacks.StreamHandler
	Executor *agents.Executor
	History  *mem.PersistentChatHistory
}

func NewAgent(config config.AppConfig) *Agent {
	agent := &Agent{}

	agent.config = config
	agent.callback = callbacks.NewStreamHandler()
	agent.LLM = config.AgentLLM()
	agent.Tools = config.GetTools()
	agent.History = mem.NewPersistentChatHistory(config)

	memoryBuffer := memory.NewConversationTokenBuffer(
		agent.LLM,
		agent.config.MemorySize(),
		memory.WithChatHistory(agent.History),
	)

	mainAgent := agents.NewConversationalAgent(
		agent.LLM,
		agent.Tools,
		agents.WithCallbacksHandler(agent.callback),
	)

	agent.Executor = agents.NewExecutor(
		mainAgent,
		agents.WithMemory(memoryBuffer),
	)

	log.Println("Agent created")
	return agent
}

func (agent *Agent) GetConvresationHistory() []string {
	agent.History.SetSession(agent.config.CurrentSession())

	history, err := agent.History.Messages(context.Background())
	if err != nil {
		log.Println("Error getting conversation history:", err)
	}

	messages := make([]string, 0)
	for _, msg := range history {
		msgContent := msg.GetContent()

		if msg.GetType() == "ai" {
			msgContent = "Clipt: " + msg.GetContent()
		}

		messages = append(messages, msgContent)
	}

	return messages
}

func (agent *Agent) Stream(ctx context.Context, callback func(ctx context.Context, chunk []byte)) {
	agent.callback.ReadFromEgress(ctx, callback)
}

func (agent *Agent) Run(ctx context.Context, input string) error {
	log.Println("Agent running with input:", input)

	agent.History.SetSession(agent.config.CurrentSession())

	if agent.Executor == nil {
		log.Println("Agent executor is nil")
	}

	_, err := chains.Run(ctx, agent.Executor, input)
	if err != nil {
		return err
	}

	return nil
}
