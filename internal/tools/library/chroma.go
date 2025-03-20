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
	client, err := chromago.NewClient("http://localhost:8000")
	if err != nil {
		log.Println("Error creating chroma client:", err)
		return nil, err
	}

	return &ChromaClient{
		client: client,
	}, nil
}

func (client *ChromaClient) InitVBD() error {

	tenant, err := client.client.GetTenant(context.Background(), "clipt-tenant")
	if err != nil {
		log.Println("Error getting tenant:", err)
	}

	if tenant == nil {
		_, err := client.client.CreateTenant(context.Background(), "clipt-tenant")
		if err != nil {
			log.Println("Error creating tenant:", err)
			return err
		}
	}

	tenantName := "clipt-tenant"
	db, err := client.client.GetDatabase(context.Background(), "clipt-vdb", &tenantName)
	if err != nil {
		log.Println("Error getting database:", err)
	}

	if db == nil {
		_, err = client.client.CreateDatabase(context.Background(), "clipt-vdb", &tenantName)
		if err != nil {
			log.Println("Error creating database:", err)
			return err
		}
	}

	client.client.SetTenant("clipt-tenant")
	client.client.SetDatabase("clipt-vdb")

	apiKey := os.Getenv("OPENAI_API_KEY")
	embeddingFunc, err := openai.NewOpenAIEmbeddingFunction(apiKey)
	if err != nil {
		log.Println("Error creating embedding function:", err)
		return err
	}

	coll, err := client.client.NewCollection(context.Background(),
		collection.WithName("documents-2"),
		collection.WithEmbeddingFunction(embeddingFunc),
		collection.WithHNSWDistanceFunction(types.L2),
		collection.WithDatabase("clipt-vdb"),
		collection.WithTenant("clipt-tenant"),
		collection.WithCreateIfNotExist(true),
	)

	if err != nil {
		log.Println("Error creating collection:", err)
		return err
	}

	client.collection = coll
	client.embeddingFunc = embeddingFunc

	return nil
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
		log.Println("Unsupported file type:", path)
		return fmt.Errorf("unsupported file type: %s", filepath.Ext(path))
	}

	recordSet, err := types.NewRecordSet(
		types.WithEmbeddingFunction(client.collection.EmbeddingFunction),
		types.WithIDGenerator(types.NewUUIDGenerator()),
	)

	if err != nil {
		log.Println("Error creating record set:", err)
		return err
	}

	// convert int64 to string
	fileSizeStr := fmt.Sprintf("%d", fileInfo.Size())

	recordSet.WithRecord(
		types.WithDocument(fileContent),
		types.WithMetadata("fileType", filepath.Ext(path)),
		types.WithMetadata("fileName", fileInfo.Name()),
		types.WithMetadata("filePath", path),
		types.WithMetadata("fileSize", fileSizeStr),
		types.WithMetadata("createdAt", fileInfo.ModTime().String()),
	)

	_, err = client.collection.AddRecords(ctx, recordSet)
	if err != nil {
		log.Println("Error adding record to collection:", err)
		return err
	}

	log.Printf("File saved to chroma DB: %s", path)
	return nil
}
