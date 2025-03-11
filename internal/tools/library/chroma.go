package library

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	openai "github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/tmc/langchaingo/documentloaders"
)

type ChromaClient struct {
	client        *chromago.Client
	embeddingFunc *openai.OpenAIEmbeddingFunction
}

func NewChromaClient() *ChromaClient {
	client, err := chromago.NewClient("http://localhost:8000",
		chromago.WithTenant("my-tenant"),
		chromago.WithDatabase("documents"),
	)

	if err != nil {
		log.Println("Error creating chroma client:", err)
		return nil
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	embeddingFunc, err := openai.NewOpenAIEmbeddingFunction(apiKey)
	if err != nil {
		log.Println("Error creating embedding function:", err)
		return nil
	}

	return &ChromaClient{
		client:        client,
		embeddingFunc: embeddingFunc,
	}
}

func (client *ChromaClient) SaveFile(ctx context.Context, path string, fileInfo os.FileInfo) error {
	var fileContent string

	switch filepath.Ext(path) {
	case ".txt":
		content, err := os.ReadFile(path)
		if err != nil {
			log.Println("Error reading file txt:", err)
			return err
		}
		fileContent = string(content)

	case ".pdf":
		fileByte, err := os.ReadFile(path)
		if err != nil {
			log.Println("Error reading file pdf:", err)
			return err
		}

		file := bytes.NewReader(fileByte)
		PDFLoader := documentloaders.NewPDF(file, file.Size())
		docs, err := PDFLoader.Load(ctx)
		fileContent = docs[0].PageContent

	}

	coll, err := client.client.NewCollection(ctx,
		collection.WithName("documents"),
		collection.WithEmbeddingFunction(client.embeddingFunc),
		collection.WithHNSWDistanceFunction(types.L2),
	)

	if err != nil {
		log.Println("Error creating collection:", err)
		return err
	}

	recordSet, err := types.NewRecordSet(
		types.WithEmbeddingFunction(client.embeddingFunc),
		types.WithIDGenerator(types.NewUUIDGenerator()),
	)

	metadata := map[string]interface{}{
		"fileType":  filepath.Ext(path),
		"fileName":  fileInfo.Name(),
		"filePath":  path,
		"fileSize":  fileInfo.Size(),
		"createdAt": fileInfo.ModTime().String(),
	}

	recordSet.WithRecord(
		types.WithDocument(fileContent),
		types.WithMetadata("metadata", metadata),
	)

	coll.AddRecords(ctx, recordSet)

	return nil
}
