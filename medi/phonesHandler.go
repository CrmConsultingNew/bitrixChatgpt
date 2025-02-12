package medi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Структура ответа от Bitrix24
type ContactResponse struct {
	Name       string `json:"NAME"`
	SecondName string `json:"SECOND_NAME"`
	LastName   string `json:"LAST_NAME"`
	Phone      []struct {
		Value string `json:"VALUE"`
	} `json:"PHONE"`
}

// ReadContactsJsonAndGetClientContactPhone читает JSON и возвращает контакты с номерами
func ReadContactsJsonAndGetClientContactPhone() []map[string]string {
	fileName := "bitrixChatgpt/birthDatesContactsLast5Days.json"
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
		return nil
	}

	var contacts []map[string]string
	if err := json.Unmarshal(data, &contacts); err != nil {
		log.Println("Ошибка при разборе JSON:", err)
		return nil
	}

	webHookUrl := "https://klinikamedi.bitrix24.ru/rest/22/m0bcq858qrmfeoa5/crm.contact.get"

	var validContacts []map[string]string

	for _, contact := range contacts {
		contactID, exists := contact["ID"]
		if !exists {
			continue
		}

		// Формируем тело запроса
		requestBody, err := json.Marshal(map[string]string{"id": contactID})
		if err != nil {
			log.Println("Ошибка при маршалинге JSON:", err)
			continue
		}

		// Создаем HTTP-запрос
		req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		// Выполняем запрос
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			continue
		}
		defer resp.Body.Close()

		// Читаем ответ
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка при чтении ответа:", err)
			continue
		}

		var response struct {
			Result ContactResponse `json:"result"`
		}

		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Println("Ошибка при разборе JSON:", err)
			continue
		}

		if len(response.Result.Phone) == 0 {
			log.Printf("Контакт ID=%s не имеет номера телефона", contactID)
			continue
		}

		rawPhone := response.Result.Phone[0].Value
		cleanPhone := formatPhoneNumber(rawPhone)
		if cleanPhone != "" {
			validContacts = append(validContacts, map[string]string{
				"phone":       cleanPhone,
				"name":        response.Result.Name,
				"second_name": response.Result.SecondName,
				"last_name":   response.Result.LastName,
			})
		}
	}

	log.Println("Обработанные контакты:", validContacts)
	return validContacts
}

// formatPhoneNumber очищает номер и приводит к корректному формату
func formatPhoneNumber(phone string) string {
	re := regexp.MustCompile(`\D`)
	cleaned := re.ReplaceAllString(phone, "")

	if len(cleaned) < 10 {
		log.Printf("Ошибка: некорректный номер %s", phone)
		return ""
	}

	if strings.HasPrefix(cleaned, "8") {
		cleaned = "7" + cleaned[1:]
	} else if !strings.HasPrefix(cleaned, "7") {
		cleaned = "7" + cleaned // Попытка добавить код России
	}

	if len(cleaned) != 11 {
		log.Printf("Ошибка: некорректный формат номера %s", phone)
		return ""
	}

	return cleaned
}
