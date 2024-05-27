/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gihutb.com/xuning888/langchain-test/z_application/chatbot/app"
	"github.com/spf13/cobra"
	"os"
)

var Model string
var MaxWindows int
var OllamaBaseUrl string

var longStr = `
This is a command-line ChatGPT based on Ollama
chatbot --model qwen:7b --maxWindows 30
chatbot --model llama3:8b --maxWindows 30 --ollamaBaseUrl http://localhost:11434
`

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:  "chatbot",
	Long: longStr,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunApplication(Model, MaxWindows, OllamaBaseUrl)
	},
}

func init() {
	// Defining flags for the command
	serveCmd.Flags().StringVar(&Model, "model", "", "Specify the model")
	serveCmd.Flags().IntVar(&MaxWindows, "maxWindows", 0, "Specify the maximum number of memory windows for the chatbot set by the user")
	serveCmd.Flags().StringVar(&OllamaBaseUrl, "ollamaBaseUrl", "http://localhost:11434", "Specify the Ollama Server url")

	// Require flags for the command
	serveCmd.MarkFlagRequired("model")
	serveCmd.MarkFlagRequired("maxWindows")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := serveCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
