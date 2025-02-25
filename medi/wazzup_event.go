package medi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	thankYouMessages = []string{
		"Спасибо за вашу оценку! Нам важно ваше мнение, поделитесь на сайте: https://medi-clinic.ru/otzyivyi/#to",
		"Благодарим вас за отзыв! Мы ценим ваше доверие. Оставьте свой отзыв здесь: https://medi-clinic.ru/otzyivyi/#to",
		"Спасибо! Ваши впечатления помогают нам становиться лучше. Поделитесь ими на сайте: https://medi-clinic.ru/otzyivyi/#to",
		"Благодарим за высокую оценку! Нам приятно, что вы довольны. Напишите отзыв: https://medi-clinic.ru/otzyivyi/#to",
		"Спасибо! Ваш отзыв важен для нас. Оставьте его на нашем сайте: https://medi-clinic.ru/otzyivyi/#to",
		"Благодарим за вашу поддержку! Мы рады стараться для вас. Поделитесь мнением: https://medi-clinic.ru/otzyivyi/#to",
		"Спасибо! Оставьте, пожалуйста, отзыв на нашем сайте: https://medi-clinic.ru/otzyivyi/#to",
		"Благодарим за оценку! Ваше мнение помогает нам расти. Напишите его здесь: https://medi-clinic.ru/otzyivyi/#to",
		"Спасибо! Нам ценно ваше доверие. Поделитесь отзывом: https://medi-clinic.ru/otzyivyi/#to",
		"Благодарим! Поделитесь вашим опытом на нашем сайте: https://medi-clinic.ru/otzyivyi/#to",
	}
	responseIndex int
	responseMutex sync.Mutex
)

type WazzupMessage struct {
	Messages []struct {
		MessageID string `json:"messageId"`
		ChatID    string `json:"chatId"`
		ChatType  string `json:"chatType"`
		Text      string `json:"text"`
		ChannelID string `json:"channelId"`
	} `json:"messages"`
}

func WazzupEventMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println("WazzupEventMessage started")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка при чтении тела запроса:", err)
		return
	}
	log.Println("Received request:", string(body))

	var event WazzupMessage
	if err := json.Unmarshal(body, &event); err != nil {
		log.Println("Ошибка при разборе JSON:", err)
		return
	}

	contacts, err := loadContacts("mediDealsAndContacts.json")
	if err != nil {
		log.Println("Ошибка загрузки контактов:", err)
		return
	}

	for _, msg := range event.Messages {
		if isValidRating(msg.Text) && isPhoneInContacts(msg.ChatID, contacts) {
			log.Println("Отправляем благодарственное сообщение пользователю:", msg.ChatID)

			responseMutex.Lock()
			messageText := thankYouMessages[responseIndex]
			responseIndex = (responseIndex + 1) % len(thankYouMessages)
			responseMutex.Unlock()

			sendMessageToWazzupReport("cf4f9e0a30ff4bb2adf92de77141c488", msg.ChannelID, msg.ChatID, msg.ChatType, messageText)
		} else {
			log.Println("isNotValid:", msg.Text)
		}
	}
}

func isValidRating(text string) bool {
	return text == "4" || text == "5"
}

func loadContacts(filename string) (map[string]string, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var contacts map[string]string
	if err := json.Unmarshal(file, &contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func isPhoneInContacts(phone string, contacts map[string]string) bool {
	for _, storedPhone := range contacts {
		if storedPhone == phone {
			return true
		}
	}
	return false
}

func sendMessageToWazzupReport(apiKey, channelId, chatId, chatType, textMessage string) {
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

	log.Printf("Ответ Wazzup: %s\n", string(body))
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Ошибка: получен статус %d\n", resp.StatusCode)
		return
	}

	log.Println("Сообщение успешно отправлено в Wazzup!")
}

//
//
//
// Всё что ниже уже не нужно (регистрация вебхука)

const (
	apiURL = "https://api.wazzup24.com/v3/webhooks"
	token  = "cf4f9e0a30ff4bb2adf92de77141c488"
)

type RequestBody struct {
	WebhooksURI   string `json:"webhooksUri"`
	Subscriptions struct {
		MessagesAndStatuses      bool `json:"messagesAndStatuses"`
		ContactsAndDealsCreation bool `json:"contactsAndDealsCreation"`
	} `json:"subscriptions"`
}

func sendPatchRequest() error {
	requestBody := RequestBody{
		WebhooksURI: "https://crmconsulting-api.ru/api/medi_wazzup_event_message",
	}
	requestBody.Subscriptions.MessagesAndStatuses = true
	requestBody.Subscriptions.ContactsAndDealsCreation = true

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	fmt.Println("PATCH request successfully sent")
	return nil
}
