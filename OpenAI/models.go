package OpenAI

type OpenAIRequest struct {
	Model    string        `json:"model"`
	Messages []MessageItem `json:"messages"`
}

type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
