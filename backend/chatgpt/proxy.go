package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ChatGPTRequest struct {
	Model     string           `json:"model"`
	Messages  []ChatGPTMessage `json:"messages"`
	MaxTokens int              `json:"max_tokens"`
}

type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Message ChatGPTMessage `json:"message"`
	} `json:"choices"`
}

// Connects to proxy and returns a configured HTTP client

func CreateProxyClient() (*http.Client, error) {
	/*proxyUser := os.Getenv("PROXY_USER")
	proxyPass := os.Getenv("PROXY_PASS")
	proxyHost := os.Getenv("PROXY_HOST")
	proxyPort := os.Getenv("PROXY_PORT")*/
	PROXY_USER := "L6QUSeH4"
	PROXY_PASS := "4fNffA4y"
	PROXY_HOST := "154.205.233.218"
	PROXY_PORT := "63415"

	proxyURL, err := url.Parse(fmt.Sprintf("socks5://%s:%s@%s:%s", PROXY_USER, PROXY_PASS, PROXY_HOST, PROXY_PORT))
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 30 * time.Second,
	}, nil
}

// Sends a question to ChatGPT API and returns the answer

func SendMessageToHuggingFace(client *http.Client, question string) (string, error) {
	url := "https://api-inference.huggingface.co/models/gpt2" // Укажите нужную модель (например, gpt2)

	apiKey := os.Getenv("HUGGING_FACE_API_KEY") // Предполагается, что API-ключ хранится в переменной окружения

	// Подготовка данных для запроса
	requestData := map[string]string{
		"inputs": question,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Декодируем ответ Hugging Face
	var response []map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Извлекаем сгенерированный текст из ответа
	if len(response) > 0 {
		if generatedText, ok := response[0]["generated_text"].(string); ok {
			return generatedText, nil
		}
	}

	return "", errors.New("no generated text in API response")
}

func SendMessageToChatGPT(client *http.Client, question string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	apiKey := os.Getenv("CHATGPT_API_KEY")

	requestData := ChatGPTRequest{
		Model: "gpt-4o",
		Messages: []ChatGPTMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: question},
		},
		MaxTokens: 50,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body)))
	}

	var response ChatGPTResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", errors.New("no content in API response")
}
