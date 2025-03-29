package callbacks

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type GraphCallbackHandler interface {
	HandleNodeStart(ctx context.Context, node string, initialState []llms.MessageContent)
	HandleNodeEnd(ctx context.Context, node string, finalState []llms.MessageContent)
	HandleNodeStream(ctx context.Context, node string, chunk []byte)
	HandleEdgeEntry(ctx context.Context, edge string, initialState []llms.MessageContent)
	HandleEdgeExit(ctx context.Context, edge string, finalState []llms.MessageContent, output string)
}

type SimpleCallbackHandler struct{}

func (callback SimpleCallbackHandler) HandleNodeStart(ctx context.Context, node string, initialState []llms.MessageContent) {
}
func (callback SimpleCallbackHandler) HandleNodeEnd(ctx context.Context, node string, finalState []llms.MessageContent) {
}
func (callback SimpleCallbackHandler) HandleNodeStream(ctx context.Context, node string, chunk []byte) {
}
func (callback SimpleCallbackHandler) HandleEdgeEntry(ctx context.Context, edge string, initialState []llms.MessageContent) {
}
func (callback SimpleCallbackHandler) HandleEdgeExit(ctx context.Context, edge string, finalState []llms.MessageContent, output string) {
}

type ReActCallbackHandler struct {
	SimpleCallbackHandler
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
