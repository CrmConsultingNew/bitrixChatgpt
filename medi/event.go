package medi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func EventHandlerMedi(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	log.Println("EventHandlerMedi was started")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Println("EventHandlerMedi raw data:", string(body))

	// Декодируем URL-кодированную строку
	decodedBody, err := url.QueryUnescape(string(body))
	if err != nil {
		log.Println("Error decoding URL:", err)
		http.Error(w, "Error decoding URL", http.StatusInternalServerError)
		return
	}

	// Парсим все параметры из строки
	values, err := url.ParseQuery(decodedBody)
	if err != nil {
		log.Println("Error parsing query:", err)
		http.Error(w, "Error parsing query", http.StatusInternalServerError)
		return
	}

	// Теперь мы можем безопасно извлечь все параметры
	event := values.Get("event")
	eventHandlerID := values.Get("event_handler_id")
	dealID := values.Get("data[FIELDS][ID]")

	log.Printf("Event: %s, EventHandlerID: %s, DealID: %s\n", event, eventHandlerID, dealID)

	// Логируем все извлеченные данные
	log.Println("Extracted values:", values)

	// Обработка dealID
	if dealID != "" {
		contactId := getDealData(dealID)
		if contactId != "" {
			contactPhone := getContactData(contactId)
			log.Println("CONTACT_PHONE is:", contactPhone)
		} else {
			log.Println("CONTACT_ID is empty")
		}
	} else {
		log.Println("dealID is empty — skipping contact lookup")
	}

	// Возвращаем ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received event: %s, handler ID: %s, deal ID: %s", event, eventHandlerID, dealID)
}

func getDealData(dealId string) (contactId string) {
	if dealId == "" {
		fmt.Println("dealId is empty")
		return ""
	}

	url := "https://klinikamedi.bitrix24.ru/rest/22/m0bcq858qrmfeoa5/crm.deal.get"

	payload := map[string]interface{}{
		"id": dealId,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %d\n", resp.StatusCode)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var result struct {
		STAGE_ID   string `json:"STAGE_ID"`
		CONTACT_ID string `json:"CONTACT_ID"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	log.Println("result.STAGE_ID:", result.STAGE_ID)

	if result.STAGE_ID == "WON" {
		log.Println("dealId is WON")
		log.Println("contactId is: ", result.CONTACT_ID)
		return result.CONTACT_ID
	} else {
		log.Println("dealId is not WON")
	}
	return ""
}

func getContactData(contactId string) string {
	if contactId == "" {
		return ""
	}

	url := "https://klinikamedi.bitrix24.ru/rest/22/m0bcq858qrmfeoa5/crm.contact.get"

	payload := map[string]interface{}{
		"id": contactId,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %d\n", resp.StatusCode)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var result struct {
		PHONE []struct {
			VALUE string `json:"VALUE"`
		} `json:"PHONE"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return ""
	}

	if len(result.PHONE) > 0 {
		return result.PHONE[0].VALUE
	}
	return "Phone not found"
}
