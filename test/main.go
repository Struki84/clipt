package main

import (
	"context"
	"fmt"
	"log"
	"os"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	openai "github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/spf13/cobra"
)

type ChromClient struct {
	client        *chroma.Client
	embeddingFunc *openai.OpenAIEmbeddingFunction
}

func NewChromaClient(apiKey string) *ChromClient {
	// Create a new Chroma client
	client, err := chroma.NewClient("http://localhost:8000")
	if err != nil {
		log.Fatalf("Error creating client: %s \n", err)
		return nil
	}

	embeddingFunc, err := openai.NewOpenAIEmbeddingFunction(apiKey, openai.WithModel("text-embedding-3-small"))
	if err != nil {
		log.Println("Error creating embedding function:", err)
		return nil
	}

	_, err = client.NewCollection(context.Background(),
		collection.WithName("test-collection"),
		collection.WithHNSWDistanceFunction(types.L2),
		collection.WithCreateIfNotExist(true),
		collection.WithEmbeddingFunction(embeddingFunc),
	)

	return &ChromClient{
		client:        client,
		embeddingFunc: embeddingFunc,
	}

}

func (chroma *ChromClient) ListCollections() {
	collections, err := chroma.client.ListCollections(context.Background())
	if err != nil {
		log.Fatalf("Error listing collections: %s \n", err)
	}

	fmt.Println("Collections:")

	for _, collection := range collections {
		fmt.Println(collection.Name)
	}
}

func (chroma *ChromClient) SaveFile(content string) {
	collection, err := chroma.client.GetCollection(context.Background(), "test-collection", chroma.embeddingFunc)
	if err != nil {
		log.Fatalf("Error creating collection: %s \n", err)
	}

	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(collection.EmbeddingFunction),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}

	rs.WithRecord(types.WithDocument(content))

	_, err = rs.BuildAndValidate(context.TODO())
	if err != nil {
		log.Fatalf("Error validating record set: %s \n", err)
	}

	// Add the records to the collection
	_, err = collection.AddRecords(context.Background(), rs)
	if err != nil {
		log.Fatalf("Error adding documents: %s \n", err)
	}
}

func (chroma *ChromClient) Search(query string) {
	collection, err := chroma.client.GetCollection(context.Background(), "test-collection", chroma.embeddingFunc)
	if err != nil {
		log.Fatalf("Error creating collection: %s \n", err)
	}

	countDocs, qrerr := collection.Count(context.TODO())
	if qrerr != nil {
		log.Fatalf("Error counting documents: %s \n", qrerr)
	}

	log.Println("countDocs: ", countDocs)

	qr, qrerr := collection.Query(context.TODO(), []string{query}, 5, nil, nil, nil)
	if qrerr != nil {
		log.Fatalf("Error querying documents: %s \n", qrerr)
	}

	fmt.Printf("qr: %v\n", qr.Documents[0][0])
}

var chromCmd = &cobra.Command{
	Use:   "chrom",
	Short: "CLI RAG tool for interacting with LLMs.",
	Long:  "",
}

func init() {

	chromaClient := NewChromaClient(os.Getenv("MAR_MAR_OPENAI_API_KEY"))

	saveFile := &cobra.Command{
		Use:   "save",
		Short: "Save a file.",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]

			chromaClient.SaveFile(input)
		},
	}

	chromCmd.AddCommand(saveFile)

	search := &cobra.Command{
		Use:   "search",
		Short: "Search for a file.",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]

			chromaClient.Search(input)
		},
	}

	chromCmd.AddCommand(search)

	list := &cobra.Command{
		Use:   "list",
		Short: "List collections.",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			chromaClient.ListCollections()
		},
	}

	chromCmd.AddCommand(list)
}

func main() {
	// if err := chromCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}

// func main() {
// 	// Create a new Chroma client
// 	client, err := chroma.NewClient("http://localhost:8000")
// 	if err != nil {
// 		log.Fatalf("Error creating client: %s \n", err)
// 		return
// 	}
//
// 	apiKey := os.Getenv("MAR_MAR_OPENAI_API_KEY")
// 	embeddingFunc, err := openai.NewOpenAIEmbeddingFunction(apiKey, openai.WithModel("text-embedding-3-small"))
// 	if err != nil {
// 		log.Println("Error creating embedding function:", err)
// 		return
// 	}
//
// 	newCollection, err := client.NewCollection(context.Background(),
// 		collection.WithName("test-collection"),
// 		collection.WithHNSWDistanceFunction(types.L2),
// 		collection.WithDatabase("clipt-vdb"),
// 		collection.WithTenant("clipt-tenant"),
// 		collection.WithCreateIfNotExist(true),
// 		collection.WithEmbeddingFunction(embeddingFunc),
// 	)
// 	if err != nil {
// 		log.Fatalf("Error creating collection: %s \n", err)
// 	}
//
// 	// Create a new record set with to hold the records to insert
// 	rs, err := types.NewRecordSet(
// 		types.WithEmbeddingFunction(newCollection.EmbeddingFunction), // we pass the embedding function from the collection
// 		types.WithIDGenerator(types.NewULIDGenerator()),
// 	)
// 	if err != nil {
// 		log.Fatalf("Error creating record set: %s \n", err)
// 	}
// 	// Add a few records to the record set
// 	rs.WithRecord(types.WithDocument("My name is John. And I have two dogs."), types.WithMetadata("key1", "value1"), types.WithID("1234"))
// 	rs.WithRecord(types.WithDocument("My name is Jane. I am a data scientist."), types.WithMetadata("key2", "value2"))
//
// 	// Build and validate the record set (this will create embeddings if not already present)
// 	_, err = rs.BuildAndValidate(context.TODO())
// 	if err != nil {
// 		log.Fatalf("Error validating record set: %s \n", err)
// 	}
//
// 	// Add the records to the collection
// 	_, err = newCollection.AddRecords(context.Background(), rs)
// 	if err != nil {
// 		log.Fatalf("Error adding documents: %s \n", err)
// 	}
//
// 	// Count the number of documents in the collection
// 	countDocs, qrerr := newCollection.Count(context.TODO())
// 	if qrerr != nil {
// 		log.Fatalf("Error counting documents: %s \n", qrerr)
// 	}
//
// 	// Query the collection
// 	fmt.Printf("countDocs: %v\n", countDocs) //this should result in 2
// 	qr, qrerr := newCollection.Query(context.TODO(), []string{"I love dogs"}, 5, nil, nil, nil)
// 	if qrerr != nil {
// 		log.Fatalf("Error querying documents: %s \n", qrerr)
// 	}
// 	fmt.Printf("qr: %v\n", qr.Documents[0][0]) //this should result in the document about dogs
// }
