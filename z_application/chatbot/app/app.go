package app

import (
	"bufio"
	"context"
	"fmt"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"log"
	"os"
	"strings"
)

var llm llms.Model
var mem *memory.ConversationWindowBuffer

func RunApplication(model string, maxWindows int, ollamaBaseUrl string) {

	if llm == nil {
		var err error
		if llm, err = ollama.New(ollama.WithModel(model),
			ollama.WithServerURL(ollamaBaseUrl)); err != nil {
			log.Fatalf("Failed to initialize Ollama client: %v", err)
		}
	}

	if mem == nil {
		mem = memory.NewConversationWindowBuffer(maxWindows)
	}

	// Create a conversation chain
	conversation := NewChineseConversation(llm, mem)
	reader := bufio.NewReader(os.Stdin)
	for {
		ctx := context.Background()
		fmt.Print(">>>")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read user input: %v\n", err)
			return
		}
		userInput = strings.TrimRight(userInput, "\n")

		if strings.HasPrefix(userInput, "/") {
			runtimeCmd, _ := strings.CutPrefix(userInput, "/")
			switch runtimeCmd {
			case "bye", "exit", "quit":
				fmt.Println("Bye bye...")
				return
			case "clearContext":
				if err2 := mem.Clear(ctx); err2 != nil {
					fmt.Println(fmt.Sprintf("Clear context error: %v", err2))
				} else {
					fmt.Println("Clear context success")
				}
				continue
			default:
				fmt.Printf("Unknown command: %s\n", runtimeCmd)
				continue
			}
		}
		if _, err2 := chains.Run(ctx, conversation, userInput,
			chains.WithStreamingFunc(makeStreamFunc()),
		); err2 != nil {
			fmt.Printf("Call llm error: %v\n", err2)
		}
		fmt.Println()
	}
}

// makeStreamFunc  make a stream callback function that prints each chunk to the console
func makeStreamFunc() func(ctx context.Context, chunk []byte) error {
	return func(ctx context.Context, chunk []byte) error {
		fmt.Printf("%s", string(chunk))
		return nil
	}
}
