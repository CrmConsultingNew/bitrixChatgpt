package medi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const filePath = "mediDealsAndContacts.json"

func EventHandlerMedi(w http.ResponseWriter, r *http.Request) {
	log.Println("EventHandlerMedi was started")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Println("EventHandlerMedi raw data:", string(body))

	decodedBody, err := url.QueryUnescape(string(body))
	if err != nil {
		log.Println("Error decoding URL:", err)
		http.Error(w, "Error decoding URL", http.StatusInternalServerError)
		return
	}

	values, err := url.ParseQuery(decodedBody)
	if err != nil {
		log.Println("Error parsing query:", err)
		http.Error(w, "Error parsing query", http.StatusInternalServerError)
		return
	}

	event := values.Get("event")
	eventHandlerID := values.Get("event_handler_id")
	dealID := values.Get("data[FIELDS][ID]")

	log.Printf("Event: %s, EventHandlerID: %s, DealID: %s\n", event, eventHandlerID, dealID)
	log.Println("Extracted values:", values)

	if dealID != "" {
		contactId := getDealData(dealID)
		if contactId != "" {
			contactPhone := getContactData(contactId)
			log.Println("CONTACT_PHONE is:", contactPhone)

			if err := updateJSONFile(dealID, contactPhone); err != nil {
				log.Println("Failed to update JSON file:", err)
			}

			message := getRandomMessage()
			log.Println("FIRST MESSAGE: >>>>>>>>> ", message)
			sendMessageToWazzupGetReport("cf4f9e0a30ff4bb2adf92de77141c488", "eec3fca0-ba9d-4bf5-89a3-35ec3080c2ae", contactPhone, "whatsapp", message)

		} else {
			log.Println("CONTACT_ID is empty")
		}
	} else {
		log.Println("dealID is empty — skipping contact lookup")
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Received event: %s, handler ID: %s, deal ID: %s", event, eventHandlerID, dealID)
}

func updateJSONFile(dealID, contactPhone string) error {
	data := make(map[string]string)

	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil && err != io.EOF {
			log.Println("Error decoding JSON file:", err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	data[dealID] = contactPhone

	file, err = os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	log.Println("JSON file updated successfully")
	return nil
}

func sendMessageToWazzupGetReport(apiKey, channelId, chatId, chatType, textMessage string) {
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

	log.Println("FULL GET_DEAL_DATA: ", string(body))

	// Корректная структура для парсинга JSON
	var response struct {
		Result struct {
			STAGE_ID   string `json:"STAGE_ID"`
			CONTACT_ID string `json:"CONTACT_ID"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	log.Println("result.STAGE_ID:", response.Result.STAGE_ID)

	if response.Result.STAGE_ID == "WON" {
		log.Println("dealId is WON")
		log.Println("contactId is: ", response.Result.CONTACT_ID)
		return response.Result.CONTACT_ID
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

	log.Println("FULL GET_CONTACT_DATA:", string(body))

	var result struct {
		Result struct {
			PHONE []struct {
				VALUE string `json:"VALUE"`
			} `json:"PHONE"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return ""
	}

	if len(result.Result.PHONE) > 0 {
		return result.Result.PHONE[0].VALUE
	}

	return "Phone not found"
}
