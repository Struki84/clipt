package library

import (
	"bytes"
	"context"
	"fmt"
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
	collection    *chromago.Collection
	embeddingFunc *openai.OpenAIEmbeddingFunction
}

func NewChromaClient() (*ChromaClient, error) {
	client, err := chromago.NewClient("http://localhost:8000",
		chromago.WithTenant("my-tenant"),
		chromago.WithDatabase("documents"),
	)

	if err != nil {
		log.Println("Error creating chroma client:", err)
		return nil, err
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	embeddingFunc, err := openai.NewOpenAIEmbeddingFunction(apiKey)
	if err != nil {
		log.Println("Error creating embedding function:", err)
		return nil, err
	}

	coll, err := client.NewCollection(context.Background(),
		collection.WithName("documents"),
		collection.WithEmbeddingFunction(embeddingFunc),
		collection.WithHNSWDistanceFunction(types.L2),
	)

	if err != nil {
		log.Println("Error creating collection:", err)
		return nil, err
	}

	return &ChromaClient{
		client:        client,
		embeddingFunc: embeddingFunc,
		collection:    coll,
	}, nil
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
		if len(docs) == 0 {
			log.Println("No content extracted from PDF:", path)
			return fmt.Errorf("empty PDF content")
		}

		fileContent = docs[0].PageContent
	default:
		log.Println("Unsupported file type:", path)
		return fmt.Errorf("unsupported file type: %s", filepath.Ext(path))
	}

	recordSet, err := types.NewRecordSet(
		types.WithEmbeddingFunction(client.embeddingFunc),
		types.WithIDGenerator(types.NewUUIDGenerator()),
	)

	if err != nil {
		log.Println("Error creating record set:", err)
		return err
	}

	metadata := map[string]interface{}{
		"fileType":  filepath.Ext(path),
		"fileName":  fileInfo.Name(),
		"filePath":  path,
		"fileSize":  fileInfo.Size(),
		"createdAt": fileInfo.ModTime().String(),
	}

	recordSet.WithRecord(
		types.WithDocument(fileContent),
		types.WithMetadatas(metadata),
	)

	_, err = client.collection.AddRecords(ctx, recordSet)
	if err != nil {
		log.Println("Error adding record to collection:", err)
		return err
	}

	return nil
}
