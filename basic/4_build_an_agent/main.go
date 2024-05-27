package main

import (
	"context"
	"gihutb.com/xuning888/langchain-test/basic/4_build_an_agent/tool"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/tools"
	"log"
)

func main() {

	calculatorAgent()
}

// testCalculator
// Note: test calculator tool
func testCalculator() {
	// Create a tool for calculator
	calculator := tools.Calculator{}
	// Prints tool name and tool description
	log.Printf("tool name: %s, tool description: %s\n", calculator.Name(), calculator.Description())
	ctx := context.Background()
	// Call calculator
	if call, err := calculator.Call(ctx, "1 + 1"); err != nil {
		log.Fatalf("call calculator error: %v", err)
		return
	} else {
		log.Printf("calculator result: %v", call)
	}
}

// calculatorAgent builds and runs a conversational agent with a calculator tool
func calculatorAgent() {
	var llm llms.Model
	var err error
	if llm, err = ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("llama3:instruct"),
	); err != nil {
		log.Fatalf("Create llm error: %v", err)
		return
	}

	mtools := []tools.Tool{
		&tool.MCalculator{},
	}

	conversationalAgent := agents.NewConversationalAgent(llm, mtools)

	executor := agents.NewExecutor(conversationalAgent, mtools)

	ctx := context.Background()

	if run, err2 := chains.Run(ctx, executor, "1 + 1 等于几？"); err2 != nil {
		log.Fatalf("call agent error: %v", err2)
		return
	} else {
		log.Printf("call agent result: %v", run)
	}
}
