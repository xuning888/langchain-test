package app

import (
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
)

const _conversationTemplate = `
以下是一段人类和AI之间的友好对话。AI非常健谈，并且会从其上下文中提供大量具体细节。如果AI不知道问题的答案，它会诚实地回答不知道。
当前对话:
{{.history}}
Human: {{.input}}
AI:`

// NewChineseConversation
// Note: This is a conversation chain, but we use Chinese to construct the prompt template.
func NewChineseConversation(llm llms.Model, memory schema.Memory) chains.LLMChain {
	return chains.LLMChain{
		Prompt: prompts.NewPromptTemplate(
			_conversationTemplate,
			[]string{"history", "input"},
		),
		LLM:          llm,
		Memory:       memory,
		OutputParser: outputparser.NewSimple(),
		OutputKey:    "text",
	}
}
