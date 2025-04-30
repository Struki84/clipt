package callbacks

import (
	"context"
	"fmt"

	"github.com/Struki84/GoLangGraph/graph"
	"github.com/tmc/langchaingo/llms"
)

type ReActCallbackHandler struct {
	graph.SimpleCallback
}

func NewReActCallbackHandler() ReActCallbackHandler {
	return ReActCallbackHandler{}
}

func (handler ReActCallbackHandler) HandleNodeStart(ctx context.Context, node string, initialState []llms.MessageContent) {
	fmt.Println(" ")
	fmt.Println("=================== " + node + " ===================")
}

func (handler ReActCallbackHandler) HandleNodeStream(ctx context.Context, node string, chunk []byte) {
	fmt.Print(string(chunk))
}
