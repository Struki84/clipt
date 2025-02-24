package library

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/textsplitter"
)

type PDFAgent struct {
	Splitter textsplitter.RecursiveCharacter
	Embedder embeddings.Embedder
	Model    llms.Model
	Docs     []models.Document
	Storage  *storage.AWS
}

func NewPDFAgent(options ...LibraryOptions) (*PDFAgent, error) {
	pdfAgent := &PDFAgent{}

	opts := LibraryAgentOptions{}
	for _, option := range options {
		option(&opts)
	}

	pdfAgent.Embedder = opts.Embedder
	pdfAgent.Model = opts.Model
	pdfAgent.Docs = opts.Docs

	pdfAgent.Splitter = textsplitter.NewRecursiveCharacter()
	pdfAgent.Splitter.ChunkSize = 500
	pdfAgent.Splitter.ChunkOverlap = 50

	pdfAgent.Storage, _ = storage.NewAWS()

	return pdfAgent, nil
}

func (agent *PDFAgent) Name() string {
	return "PDF tool."
}

func (agent *PDFAgent) Description() string {
	str := `Enables you to read PDF documents. 
	
	The tool exepects input in JSON format with search query and filename.

	Example:
	{
		"query": "Where did Simun work in 2015?",
		"file": "SimunStrukanCV.pdf"
	}

	Available documents:
	
	`
	var docStr string
	for _, doc := range agent.Docs {
		docStr += "- " + doc.Name + "\n"
	}

	str += docStr

	return str
}

func (agent *PDFAgent) Call(ctx context.Context, input string) (string, error) {
	log.Printf("PDF Agent running with input: %s", input)

	var toolInput struct {
		File  string `json:"file,omitempty"`
		Query string `json:"query,omitempty"`
	}

	re := regexp.MustCompile(`(?s)\{.*\}`)
	jsonString := re.FindString(input)

	err := json.Unmarshal([]byte(jsonString), &toolInput)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("%v: %s", "invalid input", err), nil
	}

	var requestedDoc models.Document
	for _, doc := range agent.Docs {
		if doc.Name == toolInput.File {
			requestedDoc = doc
			break
		}
	}

	fileByte, err := agent.Storage.GetFile(requestedDoc.Path)
	if err != nil {
		log.Printf("Error getting file: %s", err)
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
