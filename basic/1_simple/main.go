package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"log"
)

func main() {
	simpleExample()
	// simpleLLMApplication()
	//simpleLLMApplicationWithPromptTemplate()
}

// simpleExample demonstrates how to create and use a large language model (LLM)
// from the Ollama to generate text content. It includes setup of the LLM
// with configuration parameters, sending a "hello" message to the model, and
// handling the streaming response by printing each chunk of the generated content
// to the console. Any errors during the process are logged and will terminate the
// function.
func simpleExample() {
	// Create a llm
	var llm llms.Model
	var err error
	if llm, err = ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("llama3:8b")); err != nil {
		log.Fatalf("create llm error: %v", err)
		return
	}

	ctx := context.Background()

	// Define a stream callback function that prints each chunk to the console
	streamCallBack := func(ctx context.Context, chunk []byte) error {
		fmt.Printf("%s", string(chunk))
		return nil
	}

	// Call llm generate
	if response, generateErr := llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeHuman, "hello"), // Send "hello"
		},
		llms.WithStreamingFunc(streamCallBack), // Register a stream callback func
	); generateErr != nil {
		log.Fatalf("call llm error: %v", generateErr)
		return
	} else {
		fmt.Println()
		log.Printf("call llm success response: %s", response)
	}
}

// simpleLLMApplication is a 1_simple z_application.
// This application will translate text form english into chinese.
func simpleLLMApplication() {
	var llm llms.Model
	var err error
	if llm, err = ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("llama3:8b")); err != nil {
		log.Fatalf("create llm error: %v", err)
		return
	}

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Translate the following from English into Chinese"),
		llms.TextParts(llms.ChatMessageTypeHuman, "Hi!"),
	}

	ctx := context.Background()
	if response, generateErr := llm.GenerateContent(ctx, messages); generateErr != nil {
		log.Fatalf("call llm error: %v\n", err)
	} else {
		log.Printf("call llm success, response: %s\n", response)
	}
}

// simpleLLMApplicationWithPromptTemplate
// Use the prompt template to translate the text into the specified language.
func simpleLLMApplicationWithPromptTemplate() {
	var llm llms.Model
	var err error
	if llm, err = ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("llama3:8b")); err != nil {
		log.Fatalf("create llm error: %v", err)
		return
	}

	// Define a system prompt template
	systemPrompt := prompts.NewSystemMessagePromptTemplate("Translate the following into {{.language}}:", []string{"language"})

	promptTemplate := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		systemPrompt,
		prompts.NewHumanMessagePromptTemplate("{{.text}}", []string{"text"}),
	})

	var result string
	if result, err = promptTemplate.Format(map[string]any{
		"language": "chinese",
		"text":     "hi",
	}); err != nil {
		log.Fatalf("create prompt error: %v", err)
		return
	}
	log.Printf("prompt: %v\n", result)
	ctx := context.Background()
	if response, generateErr := llm.Call(ctx, result); generateErr != nil {
		log.Fatalf("call llm error: %v\n", err)
	} else {
		log.Printf("call llm success, response: %s\n", response)
	}
}
