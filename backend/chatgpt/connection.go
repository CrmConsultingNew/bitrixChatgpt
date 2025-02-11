package chatgpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
)

func SendRequest(vueReq string) string {
	apiKey := os.Getenv("CHATGPT_API_KEY")
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: vueReq,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}
	var text = resp.Choices[0].Message.Content
	fmt.Println(resp.Choices[0].Message.Content)
	return text
}
