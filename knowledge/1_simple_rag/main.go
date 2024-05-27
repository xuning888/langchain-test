package main

import (
	"context"
	"fmt"
	chromago "github.com/amikos-tech/chroma-go"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/chroma"
	"log"
	"os"
)

var embeddingClient embeddings.EmbedderClient

var llm llms.Model

func init() {
	var err error
	if embeddingClient, err = ollama.New(
		ollama.WithModel("mxbai-embed-large:latest"),
	); err != nil {
		log.Fatalln(fmt.Sprintf("create embedding client error: %v", err))
		return
	}
	if llm, err = ollama.New(ollama.WithModel("qwen:7b")); err != nil {
		log.Fatalln(fmt.Sprintf("create llm error: %v", err))
		return
	}
}

var ragPrompt = `You are an assistant for question-answering tasks. Use the following pieces of retrieved context to answer the question. If you don't know the answer, just say that you don't know. Use three sentences maximum and keep the answer concise.
Question: {{.question}}
Context: {{.context}}
Answer:`

func main() {
	// TODO
	// https://lilianweng.github.io/posts/2023-06-23-agent/
	var file *os.File
	var err error
	if file, err = os.Open("/Users/xuning/gopath/src/gihutb.com/xuning888/langchain-test/knowledge/1_simple_rag/testdata/Redis源代码分析.pdf"); err != nil {
		log.Fatalln(fmt.Sprintf("load file error: %v", err))
		return
	}

	stat, _ := os.Stat("/Users/xuning/gopath/src/gihutb.com/xuning888/langchain-test/knowledge/1_simple_rag/testdata/Redis源代码分析.pdf")

	pdfLoader := documentloaders.NewPDF(file, stat.Size())
	var documents []schema.Document
	if documents, err = pdfLoader.Load(context.Background()); err != nil {
		log.Fatalln(fmt.Sprintf("load html error: %v", err))
		return
	}
	textSplitter := textsplitter.NewRecursiveCharacter(textsplitter.WithChunkSize(1000), textsplitter.WithChunkOverlap(200))

	var splitDocuments []schema.Document
	if splitDocuments, err = textsplitter.SplitDocuments(textSplitter, documents); err != nil {
		log.Fatalln(fmt.Sprintf("SplitDocuments error: %v", err))
		return
	}

	var embedder embeddings.Embedder
	if embedder, err = embeddings.NewEmbedder(embeddingClient); err != nil {
		log.Fatalln(fmt.Sprintf("Create embedder error: %v", err))
		return
	}

	namespace := uuid.New().String()

	chromaClient, _ := chromago.NewClient("http://localhost:8000")

	defer func() {
		_, deleteCollectionErr := chromaClient.DeleteCollection(context.Background(), namespace)
		if deleteCollectionErr != nil {
			log.Printf("delete collection %s, error: %v", namespace, deleteCollectionErr)
		}
	}()

	store, err := chroma.New(
		chroma.WithChromaURL("http://localhost:8000"),
		chroma.WithEmbedder(embedder),
		chroma.WithDistanceFunction("cosine"),
		chroma.WithNameSpace(namespace),
	)

	ctx := context.Background()
	if _, err2 := store.AddDocuments(ctx, splitDocuments); err2 != nil {
		log.Fatalln(fmt.Sprintf("AddDocuments error: %v", err2))
		return
	}

	//chains.NewRetrievalQA()

	question := "What is Task Decomposition?"
	search, err := store.SimilaritySearch(ctx, question, 1)
	if err != nil {
		log.Fatalln(fmt.Sprintf("search error: %v", err))
		return
	}

	template := prompts.NewPromptTemplate(ragPrompt, []string{"question", "context"})
	chain := chains.NewLLMChain(llm, template)

	call, err := chains.Call(ctx, chain, map[string]any{
		"question": question,
		"context":  search[0].PageContent,
	})
	if err != nil {
		log.Fatalln(fmt.Sprintf("call llm error: %v", err))
		return
	}

	fmt.Printf("%s", call)
}
