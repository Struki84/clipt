package callbacks

import (
	"context"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/schema"
)

// DefaultKeywords is map of the agents final out prefix keywords.
//
//nolint:all
var DefaultKeywords = []string{"Final Answer:", "Final:", "AI:"}

// nolint:all
type StreamHandler struct {
	callbacks.SimpleHandler
	egress          chan []byte
	Keywords        []string
	LastTokens      string
	KeywordDetected bool
	PrintOutput     bool

	//tmp fix
	ChainsActive   []string
	ChainsFinished []string
}

var _ callbacks.Handler = &StreamHandler{}

func NewStreamHandler(keywords ...string) *StreamHandler {
	file, err := tea.LogToFile("./debug.log", "debug")
	if err != nil {
		log.Println("Streamer handler error:", err)
	}

	defer file.Close()

	if len(keywords) > 0 {
		DefaultKeywords = keywords
	}

	return &StreamHandler{
		egress:         make(chan []byte),
		Keywords:       DefaultKeywords,
		ChainsActive:   make([]string, 0),
		ChainsFinished: make([]string, 0),
	}
}

func (handler *StreamHandler) GetEgress() chan []byte {
	return handler.egress
}

func (handler *StreamHandler) ReadFromEgress(ctx context.Context, callback func(ctx context.Context, chunk []byte)) {
	go func() {
		defer close(handler.egress)
		for data := range handler.egress {
			callback(ctx, data)
		}
	}()
}

func (handler *StreamHandler) HandleChainStart(ctx context.Context, inputs map[string]any) {
	handler.PrintOutput = false
	handler.KeywordDetected = false

	log.Println("Chain started")

}

func (handler *StreamHandler) HandleChainEnd(ctx context.Context, outputs map[string]any) {

	log.Println("Chain finished")

}

func (handler *StreamHandler) HandleChainError(ctx context.Context, err error) {

	log.Println("Chain error:", err)

}

func (handler *StreamHandler) HandleAgentAction(ctx context.Context, action schema.AgentAction) {

	log.Println("Agent action")

}

func (handler *StreamHandler) HandleAgentFinish(ctx context.Context, finish schema.AgentFinish) {

	log.Println("Agent finished")

}

func (handler *StreamHandler) HandleStreamingFunc(ctx context.Context, chunk []byte) {
	chunkStr := string(chunk)
	handler.LastTokens += chunkStr

	// Buffer the last few chunks to match the longest keyword size
	longestSize := len(handler.Keywords[0])
	for _, k := range handler.Keywords {
		if len(k) > longestSize {
			longestSize = len(k)
		}
	}

	if len(handler.LastTokens) > longestSize {
		handler.LastTokens = handler.LastTokens[len(handler.LastTokens)-longestSize:]
	}

	// Check for keywords
	for _, k := range DefaultKeywords {
		if strings.Contains(handler.LastTokens, k) {
			handler.KeywordDetected = true
		}
	}

	// Check for colon and set print mode.
	if handler.KeywordDetected && chunkStr != ":" {
		handler.PrintOutput = true
	}

	// Print the final output after the detection of keyword.
	if handler.PrintOutput {
		handler.egress <- chunk
	}
}
