package medi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

var GlobalTextMessageToClient string

func SendMessageToClient(phoneNumber string, message string) {
	//sendMessageToSMSRU("9A57992F-5DF0-9F20-BC3C-7C38EF743CA7", phoneNumber, message)
	//SendMessageToWazzup("cf4f9e0a30ff4bb2adf92de77141c488", "eec3fca0-ba9d-4bf5-89a3-35ec3080c2ae", phoneNumber, "whatsapp", message)
}

// SendMessageToWazzup отправляет сообщение через Wazzup API
func SendMessageToWazzup(apiKey, channelId, chatId, chatType, textMessage string) {
	log.Println("sendmessageToWazzup was started....")
	url := "https://api.wazzup24.com/v3/message"

	requestBody, err := json.Marshal(map[string]interface{}{
		"channelId": channelId,
		"chatId":    chatId,
		"chatType":  chatType,
		"text":      textMessage,
	})
	if err != nil {
		log.Println("Ошибка при маршалинге JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Ошибка при создании HTTP-запроса:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Логируем полный ответ от Wazzup
	log.Printf("Ответ Wazzup: %s\n", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Ошибка: получен статус %d\n", resp.StatusCode)
		return
	}

	log.Println("Сообщение успешно отправлено в Wazzup!")
}

// sendMessageToSMSRU отправляет SMS через SMS.RU API
func sendMessageToSMSRU(apiID, phoneNumber, textMessage string) {
	log.Println("sendMessageToSMSRU was started....")

	baseURL := "https://sms.ru/sms/send"

	// Формируем параметры запроса
	params := url.Values{}
	params.Set("api_id", apiID)
	params.Set("to", phoneNumber)
	params.Set("msg", textMessage)
	params.Set("json", "1")

	// Выполняем POST-запрос
	resp, err := http.PostForm(baseURL, params)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Логируем ответ
	log.Printf("Ответ SMS.RU: %s\n", string(body))

	// Проверяем статус HTTP
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка: получен статус %d\n", resp.StatusCode)
		return
	}

	log.Println("Сообщение успешно отправлено через SMS.RU!")
}
