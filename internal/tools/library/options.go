package library

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
)

type LibraryAgentOptions struct {
	Embedder embeddings.Embedder
	Model    llms.Model
}

type LibraryOptions func(*LibraryAgentOptions)

func WithEmbedder(embedder embeddings.Embedder) LibraryOptions {
	return func(options *LibraryAgentOptions) {
		options.Embedder = embedder
	}
}

func WithModel(model llms.Model) LibraryOptions {
	return func(options *LibraryAgentOptions) {
		options.Model = model
	}
}
