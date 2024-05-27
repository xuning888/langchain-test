package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"log"
)

var llm llms.Model

var embeddingClient embeddings.EmbedderClient

func init() {
	if Ollama, err := ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("qwen:7b")); err != nil {
		log.Fatalf("create llm error: %v", err)
		return
	} else {
		llm = Ollama
		embeddingClient = Ollama
	}
}

func main() {
	// Embedding()
	// EmbeddingAndSimilaritySearch()
}

func Embedding() {
	ctx := context.Background()
	if embedding, err := embeddingClient.CreateEmbedding(ctx, []string{"hello"}); err != nil {
		log.Fatalf("embedding error: %v", err)
		return
	} else {
		fmt.Println(fmt.Sprintf("result: %v", embedding))
	}
}

// EmbeddingAndSimilaritySearch
// Note: embedding the document and test SimilaritySearch
func EmbeddingAndSimilaritySearch() {
	// Prepare documents
	documents := []schema.Document{
		{
			PageContent: "Dogs are great companions, known for their loyalty and friendliness.",
			Metadata:    map[string]interface{}{"source": "mammal-pets-doc"},
		},
		{
			PageContent: "Cats are independent pets that often enjoy their own space.",
			Metadata:    map[string]interface{}{"source": "mammal-pets-doc"},
		},
		{
			PageContent: "Goldfish are popular pets for beginners, requiring relatively simple care.",
			Metadata:    map[string]interface{}{"source": "fish-pets-doc"},
		},
		{
			PageContent: "Parrots are intelligent birds capable of mimicking human speech.",
			Metadata:    map[string]interface{}{"source": "bird-pets-doc"},
		},
		{
			PageContent: "Rabbits are social animals that need plenty of space to hop around.",
			Metadata:    map[string]interface{}{"source": "mammal-pets-doc"},
		},
	}

	// Create a embedder
	embedder, err := embeddings.NewEmbedder(embeddingClient)
	if err != nil {
		log.Fatalf("create embedder error: %v\n", err)
		return
	}

	// Define a random string as namespace
	namespace := uuid.New().String()

	// Create a chroma vector db
	store, err := chroma.New(
		chroma.WithChromaURL("http://localhost:8000"),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace(namespace),
	)

	ctx := context.Background()
	// AddDocuments
	docIds, err := store.AddDocuments(ctx, documents)
	if err != nil {
		log.Fatalf("AddDocuments error: %v", err)
		return
	}

	log.Println("AddDocuments success")
	for _, docId := range docIds {
		fmt.Println(fmt.Sprintf("docId: %v", docId))
	}

	// SimilaritySearch
	docs, err := store.SimilaritySearch(ctx, "Cats", 3)
	if err != nil {
		log.Fatalf("SimilaritySearch error: %v", err)
		return
	}

	for _, doc := range docs {
		fmt.Println(doc)
	}
}
