package medi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func WazzupEventMessage(w http.ResponseWriter, r *http.Request) {
	// Отправляем ответ 200 OK сразу, чтобы API Wazzup не посчитал запрос невалидным
	w.WriteHeader(http.StatusOK)

	log.Println("WazzupEventMessage was startedz <-")
	// Запускаем отправку PATCH-запроса в фоне
	go func() {
		err := sendPatchRequest()
		if err != nil {
			log.Println("Ошибка при отправке PATCH запроса:", err)
		}
	}()

	// Читаем тело запроса и логируем его
	rdr, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка при чтении тела запроса:", err)
		return
	}
	log.Println("WazzupEventMessage:", string(rdr))
}

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
