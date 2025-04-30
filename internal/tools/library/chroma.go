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
	embeddingFunc *openai.OpenAIEmbeddingFunction
}

func NewChromaClient() (*ChromaClient, error) {
	client, err := chromago.NewClient("http://localhost:8000")
	if err != nil {
		log.Println("Error creating chroma client:", err)
		return nil, err
	}

	embeddingFunc, err := openai.NewOpenAIEmbeddingFunction(os.Getenv("MAR_MAR_OPENAI_API_KEY"), openai.WithModel("text-embedding-3-large"))
	if err != nil {
		log.Println("Error creating embedding function:", err)
		return nil, err
	}

	_, err = client.NewCollection(context.Background(),
		collection.WithName("clipt_documents"),
		collection.WithHNSWDistanceFunction(types.L2),
		collection.WithCreateIfNotExist(true),
		collection.WithEmbeddingFunction(embeddingFunc),
	)

	return &ChromaClient{
		client:        client,
		embeddingFunc: embeddingFunc,
	}, nil
}

func (client *ChromaClient) SaveFile(ctx context.Context, path string, fileInfo os.FileInfo) error {
	log.Println("Saving file to chroma DB:", path)
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
		if err != nil {
			log.Println("Error loading PDF:", err)
			return err
		}

		if len(docs) == 0 {
			log.Println("No content extracted from PDF:", path)
			return fmt.Errorf("empty PDF content")
		}

		fileContent = docs[0].PageContent
	default:
		return fmt.Errorf("unsupported file type: %s", filepath.Ext(path))
	}

	collection, err := client.client.GetCollection(ctx, "clipt_documents", client.embeddingFunc)
	if err != nil {
		log.Println("Error getting collection:", err)
		return err
	}

	recordSet, err := types.NewRecordSet(
		types.WithEmbeddingFunction(collection.EmbeddingFunction),
		types.WithIDGenerator(types.NewUUIDGenerator()),
	)

	if err != nil {
		log.Println("Error creating record set:", err)
		return err
	}

	// convert int64 to string
	// fileSizeStr := fmt.Sprintf("%d", fileInfo.Size())

	recordSet.WithRecord(
		types.WithDocument(fileContent),
	)

	_, err = recordSet.BuildAndValidate(ctx)
	if err != nil {
		log.Println("Error validating record set:", err)
		return err
	}

	_, err = collection.AddRecords(ctx, recordSet)
	if err != nil {
		return err
	}

	log.Printf("File saved to chroma DB: %s", path)
	return nil
}

func (client *ChromaClient) GetFile(ctx context.Context, ID string) (string, error) {
	log.Println("Getting file from chroma DB:", ID)

	return "", nil
}

func (client *ChromaClient) DeleteFile(ctx context.Context, ID string) error {
	log.Println("Deleting file from chroma DB:", ID)
	return nil
}

func (client *ChromaClient) SearchFiles(ctx context.Context, query string) ([]string, error) {
	log.Println("Searching for files in chroma DB:", query)
	return nil, nil
}
