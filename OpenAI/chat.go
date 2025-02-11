package OpenAI

import (
	"bitrix_app/backend/chatgpt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ChatGptTest(prompt, message string) (string, error) {
	// Создаём прокси-клиент
	log.Println("openai started....")
	proxyClient, err := chatgpt.CreateProxyClient()
	if err != nil {
		log.Println("Ошибка создания HTTP клиента с прокси:", err)
		return "", err
	}

	// Формируем запрос к OpenAI API
	requestPayload := OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []MessageItem{
			{
				Role:    "system",
				Content: prompt, // Используем переданный промпт
			},
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		log.Println("Ошибка при сериализации запроса:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", OpenAIAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Ошибка при создании HTTP-запроса:", err)
		return "", err
	}
	openApiKey := os.Getenv("CHATGPT_API_KEY")
	req.Header.Set("Authorization", "Bearer "+openApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос через прокси-клиент
	resp, err := proxyClient.Do(req)
	if err != nil {
		log.Println("Ошибка при отправке запроса в OpenAI API:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Читаем ответ
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении ответа от OpenAI API:", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка от OpenAI API: %s\n", string(responseData))
		return "", err
	}

	// Разбираем ответ
	var openAIResponse OpenAIResponse
	if err := json.Unmarshal(responseData, &openAIResponse); err != nil {
		log.Println("Ошибка при разборе JSON-ответа:", err)
		return "", err
	}

	var answer string
	// Проверяем и извлекаем ответ
	if len(openAIResponse.Choices) > 0 {
		answer = openAIResponse.Choices[0].Message.Content
	} else {
		log.Println("Ответ от ChatGPT пуст.")
	}
	return answer, nil
}
