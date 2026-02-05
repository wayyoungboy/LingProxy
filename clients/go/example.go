package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lingproxy/lingproxy/clients/go/client"
	"github.com/openai/openai-go/v3"
)

func main() {
	// Initialize client
	// Option 1: Use environment variable
	// export LINGPROXY_API_KEY=your-api-key
	// c, err := client.NewClient(&client.ClientOptions{})

	// Option 2: Pass API key directly
	apiKey := os.Getenv("LINGPROXY_API_KEY")
	if apiKey == "" {
		apiKey = "ling-Uc9tFvNr97HaMXrKXw2R1ZqNiHt_pp0M_OsDOvjns8M="
		fmt.Println("Warning: Using default API key. Please set LINGPROXY_API_KEY environment variable for production use.")
		fmt.Println()
	}

	c, err := client.NewClient(&client.ClientOptions{
		APIKey:  apiKey,
		BaseURL: "http://localhost:8080/llm/v1",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	ctx := context.Background()

	fmt.Println("=" + repeat("=", 60))
	fmt.Println("LingProxy Go Client Demo")
	fmt.Println("=" + repeat("=", 60))
	fmt.Println()

	// Example 1: List models
	fmt.Println("Example 1: List available models")
	fmt.Println(repeat("-", 70))
	models, err := c.ListModels(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d models:\n", len(models))
		for i, model := range models {
			if i >= 5 {
				break
			}
			fmt.Printf("  - %s\n", model.ID)
		}
	}
	fmt.Println()

	// Example 2: Chat completion
	fmt.Println("Example 2: Chat completion")
	fmt.Println(repeat("-", 70))
	temp := 0.7
	maxTokens := int64(100)
	chatReq := client.ChatCompletionRequest{
		Model: "glm-4.5-flash", // Replace with your model name
		Messages: []client.ChatMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Tell me about Greece's largest city."},
		},
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	}

	chatResp, err := c.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Response ID: %s\n", chatResp.ID)
		fmt.Printf("Model: %s\n", chatResp.Model)
		if len(chatResp.Choices) > 0 {
			content := chatResp.Choices[0].Message.Content
			if len(content) > 200 {
				fmt.Printf("Content: %s...\n", content[:200])
			} else {
				fmt.Printf("Content: %s\n", content)
			}
		}
		fmt.Printf("Usage: Prompt=%d, Completion=%d, Total=%d\n",
			chatResp.Usage.PromptTokens,
			chatResp.Usage.CompletionTokens,
			chatResp.Usage.TotalTokens)
	}
	fmt.Println()

	// Example 3: Using OpenAI SDK style (direct access)
	fmt.Println("Example 3: Using OpenAI SDK style")
	fmt.Println(repeat("-", 70))
	openaiClient := c.GetClient()
	response, err := openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "glm-4.5-flash",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Hello! How are you?"),
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		if len(response.Choices) > 0 {
			fmt.Printf("Response: %s\n", response.Choices[0].Message.Content)
		}
	}
	fmt.Println()

	fmt.Println("=" + repeat("=", 60))
	fmt.Println("Demo completed!")
	fmt.Println("=" + repeat("=", 60))
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
