package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"log"
)

var llm llms.Model

func init() {
	var err error
	if llm, err = ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("qwen:7b")); err != nil {
		log.Fatalf("create llm error: %v", err)
		return
	}
}

func main() {
	// demonstrateAIMemoryLimitations()
	// exampleWithMemory()
	// exampleWithCustomMessageHistory()
	exampleConversation()
}

// exampleConversation
// 使用langchain提供的 memory 组件构建一个聊天的chain
func exampleConversation() {
	mem := memory.NewConversationWindowBuffer(10)

	conversation := chains.NewConversation(llm, mem)
	ctx := context.Background()
	if run, err := chains.Run(ctx, conversation, "Hi! my name is 老铁"); err != nil {
		log.Fatalf("call llm error: %v", err)
		return
	} else {
		fmt.Println(run)
	}

	if run, err := chains.Run(ctx, conversation, "我的名字是什么？"); err != nil {
		log.Fatalf("call llm error: %v", err)
		return
	} else {
		fmt.Println(run)
	}
}

// demonstrateAIMemoryLimitations
// 这个例子说明了AI没有记忆
func demonstrateAIMemoryLimitations() {
	ctx := context.Background()
	if response, err := llm.Call(ctx, "Hi! my name is 老铁"); err != nil {
		log.Fatalf("call llm error: %v", err)
		return
	} else {
		log.Println(response)
	}
	if response, err := llm.Call(ctx, "我的名字是什么？"); err != nil {
		log.Fatalf("call llm error: %v", err)
		return
	} else {
		log.Println(response)
	}
}

// exampleWithMemory
// 这个例子说明了如果需要AI有记忆就需要在上下文中告诉它
func exampleWithMemory() {
	ctx := context.Background()
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, "Hi! my name is 老铁"),
		llms.TextParts(llms.ChatMessageTypeAI, "Hello 老铁! How can I assist you today?"),
		llms.TextParts(llms.ChatMessageTypeHuman, "我的名字是什么?"),
	}
	if response, err := llm.GenerateContent(ctx, messages); err != nil {
		log.Fatalf("call llm error: %v", err)
		return
	} else {
		log.Printf("%s", response)
	}
}

// exampleWithCustomMessageHistory
// 基于 exampleWithMemory 我们自己维护一个MessageHistory告诉llm聊天的上下文
func exampleWithCustomMessageHistory() {
	messageHistory := &MessageHistory{
		llm:     llm,
		history: make([]llms.MessageContent, 0),
	}
	ctx := context.Background()
	if call, err := messageHistory.Call(ctx, "Hi! my name is 老铁"); err != nil {
		log.Fatalf("call llm error: %v", err)
		return
	} else {
		fmt.Println(call)
	}
	if call, err := messageHistory.Call(ctx, "我的名字是什么？"); err != nil {
		log.Fatalf("call llm error: %v", err)
	} else {
		fmt.Println(call)
	}
}

type MessageHistory struct {
	llm     llms.Model
	history []llms.MessageContent
}

func (m *MessageHistory) Call(ctx context.Context, input string) (string, error) {
	history := m.GetHistory()
	humanMessage := llms.TextParts(llms.ChatMessageTypeHuman, input)
	messages := append(history, humanMessage)
	response, err := m.llm.GenerateContent(ctx, messages)
	if err != nil {
		return "", err
	}
	m.history = append(m.history, humanMessage, llms.TextParts(llms.ChatMessageTypeAI, response.Choices[0].Content))
	return response.Choices[0].Content, nil
}

func (m *MessageHistory) GetHistory() []llms.MessageContent {
	return m.history
}
