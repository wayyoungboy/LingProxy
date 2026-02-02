package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	// 配置参数
	// ⚠️ 请从环境变量或配置文件读取 API Key，不要硬编码
	apiKey := os.Getenv("LINGPROXY_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: LINGPROXY_API_KEY environment variable is not set")
		fmt.Println("Please set it with: export LINGPROXY_API_KEY=your_api_key")
		os.Exit(1)
	}
	baseURL := "http://localhost:8080/llm/v1/"
	model := "glm-4.5-flash"

	// 创建 OpenAI 客户端
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)
	ctx := context.Background()

	// 测试聊天补全
	testChatCompletion(ctx, client, model)
}

// testChatCompletion 测试聊天补全功能
func testChatCompletion(ctx context.Context, client openai.Client, model string) {
	fmt.Println("Testing Chat Completion...")
	fmt.Println("Model:", model)
	fmt.Println()

	// 创建聊天消息
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant."),
		openai.UserMessage("Tell me about Greece's largest city."),
	}

	// 创建聊天补全请求
	params := openai.ChatCompletionNewParams{
		Model:    model,
		Messages: messages,
	}

	// 发送请求
	response, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		fmt.Printf("Error creating chat completion: %v\n", err)
		return
	}

	// 打印响应
	fmt.Println("Response:")
	fmt.Println("ID:", response.ID)
	fmt.Println("Created:", response.Created)
	fmt.Println("Model:", response.Model)
	fmt.Println("Content:", response.Choices[0].Message.Content)
	fmt.Println()

	// 打印使用情况
	fmt.Println("Usage:")
	fmt.Println("Prompt Tokens:", response.Usage.PromptTokens)
	fmt.Println("Completion Tokens:", response.Usage.CompletionTokens)
	fmt.Println("Total Tokens:", response.Usage.TotalTokens)
	fmt.Println()
}
