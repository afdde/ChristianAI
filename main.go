package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	// Define available models
	models := []string{
		"meta-llama/llama-4-maverick",
		"google/gemini-2.5-pro-preview-03-25",
		"google/gemini-2.5-flash-preview:thinking",
		"openai/gpt-4o",
		"anthropic/claude-3.7-sonnet:beta",
	}

	// Display available models
	fmt.Println("Available models:")
	for i, model := range models {
		fmt.Printf("%d. %s\n", i, model)
	}

	// Get model selection from user
	var modelIndex int64
	for {
		fmt.Print("Select a model ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		input = strings.TrimSpace(input) // Remove both newline and spaces

		// Convert string to int64
		num, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			fmt.Println("Please enter a valid number")
			continue
		}

		modelIndex = num
		if modelIndex < 0 || modelIndex > int64(len(models)) {
			fmt.Printf("fuck you %d\n", len(models))
			continue
		}
		break
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your  message: ")
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message) // Remove both newline and spaces

	// Read system message from prompt.txt
	promptBytes, err := os.ReadFile("prompt.txt")
	if err != nil {
		fmt.Printf("Error reading prompt.txt: %v\n", err)
		return
	}
	systemMessage := string(promptBytes)

	client := openai.NewClient(
		option.WithAPIKey(""), // defaults to os.LookupEnv("OPENAI_API_KEY")
		option.WithBaseURL("https://openrouter.ai/api/v1"),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
			openai.SystemMessage(systemMessage),
		},
		Model: models[modelIndex],
	})
	if err != nil {
		fmt.Printf("Error details: %+v\n", err)
		return
	}

	if len(chatCompletion.Choices) == 0 {
		fmt.Println("No choices returned from API")
		return
	}
	println(chatCompletion.Choices[0].Message.Content)
}
