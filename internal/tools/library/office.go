package library

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/struki84/clipt/internal/loaders"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/textsplitter"
)

type OfficeAgent struct {
	Splitter textsplitter.RecursiveCharacter
	Embedder embeddings.Embedder
	Model    llms.Model
}

func NewOfficeTool(options ...LibraryOptions) (*OfficeAgent, error) {
	opts := LibraryAgentOptions{}
	for _, option := range options {
		option(&opts)
	}

	officeAgent := &OfficeAgent{}

	officeAgent.Embedder = opts.Embedder
	officeAgent.Model = opts.Model
	officeAgent.Splitter = textsplitter.NewRecursiveCharacter()
	officeAgent.Splitter.ChunkSize = 500
	officeAgent.Splitter.ChunkOverlap = 50

	return officeAgent, nil
}

func (agent *OfficeAgent) Name() string {
	return "Office tool."
}

func (agent *OfficeAgent) Description() string {
	str := `Enables you to read Office documents.

	The tool exepects input in JSON format with search query and filename.

	Example:
	{
		"query": "Where did Simun work in 2015?",
		"file": "SimunStrukanCV.docx"
	}
	`
	return str
}

func (agent *OfficeAgent) Call(ctx context.Context, input string) (string, error) {
	log.Printf("Office Agent running with input: %s", input)

	var toolInput struct {
		File  string `json:"file"`
		Query string `json:"query"`
	}

	re := regexp.MustCompile(`(?s)\{.*\}`)
	jsonString := re.FindString(input)

	err := json.Unmarshal([]byte(jsonString), &toolInput)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s", err)
		return fmt.Sprintf("%v: %s", "invalid input", err), nil
	}

	dirPath := "./files"

	fileByte, err := os.ReadFile(dirPath + "/" + toolInput.File)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return "", err
	}

	file := bytes.NewReader(fileByte)
	officeLoader := loaders.NewOffice(file, file.Size(), toolInput.File)

	docs, err := officeLoader.LoadAndSplit(ctx, agent.Splitter)
	if err != nil {
		log.Printf("Error loading and splitting: %s", err)
		return "", err
	}

	QAChain := chains.LoadStuffQA(agent.Model)

	answer, err := chains.Call(ctx, QAChain, map[string]any{
		"input_documents": docs,
		"question":        toolInput.Query,
	})
	if err != nil {
		log.Printf("Error running QAChain: %s", err)
		return "", err
	}

	response := answer["text"].(string)
	return response, nil
}
