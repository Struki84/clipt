package library

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/textsplitter"
)

type PDFReaderTool struct {
	Splitter textsplitter.RecursiveCharacter
	Embedder embeddings.Embedder
	Model    llms.Model
}

func NewPDFReaderTool(options ...LibraryOptions) (*PDFReaderTool, error) {
	pdfAgent := &PDFReaderTool{}

	opts := LibraryAgentOptions{}
	for _, option := range options {
		option(&opts)
	}

	pdfAgent.Embedder = opts.Embedder
	pdfAgent.Model = opts.Model

	pdfAgent.Splitter = textsplitter.NewRecursiveCharacter()
	pdfAgent.Splitter.ChunkSize = 500
	pdfAgent.Splitter.ChunkOverlap = 50

	return pdfAgent, nil
}

func (agent *PDFReaderTool) Name() string {
	return "PDF tool."
}

func (agent *PDFReaderTool) Description() string {
	str := `Enables you to read PDF documents.

	The tool exepects input in JSON format with search query and filename.

	Example:
	{
		"query": "Where did Simun work in 2015?",
		"file": "SimunStrukanCV.pdf"
	}
	`

	return str
}

func (agent *PDFReaderTool) Call(ctx context.Context, input string) (string, error) {
	log.Printf("PDF reader tool running with input: %s", input)

	var toolInput struct {
		File  string `json:"file,omitempty"`
		Query string `json:"query,omitempty,omitempty"`
	}

	err := json.Unmarshal([]byte(input), &toolInput)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("%v: %s", "invalid input", err), nil
	}

	dirPath := "./files"
	fileByte, err := os.ReadFile(dirPath + "/" + toolInput.File)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return "", err
	}

	file := bytes.NewReader(fileByte)
	PDFLoader := documentloaders.NewPDF(file, file.Size())
	docs, err := PDFLoader.LoadAndSplit(ctx, agent.Splitter)
	if err != nil {
		log.Printf("Error loading and splitting: %s", err)
		return "", err
	}

	QAChain := chains.LoadStuffQA(agent.Model)

	if toolInput.Query == "" {
		toolInput.Query = "Provide the a summary of the document"
	}

	answer, err := chains.Call(ctx, QAChain, map[string]any{
		"input_documents": docs,
		"question":        toolInput.Query,
	})
	if err != nil {
		log.Printf("Error calling chain: %s", err)
		return "", err
	}

	response := answer["text"].(string)
	return response, nil
}
