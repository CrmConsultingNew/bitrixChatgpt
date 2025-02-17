package medi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

	// Парсим форму
	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Правильно извлекаем dealID
	dealID := r.Form.Get("data[FIELDS][ID]")
	event := r.Form.Get("event")
	eventHandlerID := r.Form.Get("event_handler_id")

	// Логируем полученные значения
	log.Printf("Event: %s, EventHandlerID: %s, DealID: %s\n", event, eventHandlerID, dealID)

	// Возвращаем ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received event: %s, handler ID: %s, deal ID: %s", event, eventHandlerID, dealID)

	// handleContact
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

	if result.STAGE_ID == "WON" {
		return result.CONTACT_ID
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
