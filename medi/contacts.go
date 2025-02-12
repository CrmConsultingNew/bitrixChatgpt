package medi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetContactsListWithCustomFieldsBirthdate() {
	webHookUrl := "https://klinikamedi.bitrix24.ru/rest/22/m0bcq858qrmfeoa5/crm.contact.list"

	// Получаем дату 5 дней назад
	pastDate := time.Now().AddDate(0, 0, +5)
	currentDay := strconv.Itoa(pastDate.Day())
	currentMonth := strconv.Itoa(int(pastDate.Month()))

	contacts := []map[string]string{}
	start := 0

	for {
		requestBody := map[string]interface{}{
			"FILTER": map[string]string{
				"=UF_CRM_1739257086952": currentDay,   // День рождения за 5 дней до сегодня
				"=UF_CRM_1739257138715": currentMonth, // Месяц дня рождения
			},
			"SELECT": []string{"ID"},
			"start":  start,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Println("Ошибка при маршалинге JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			return
		}
		defer resp.Body.Close()

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка при чтении ответа:", err)
			return
		}

		log.Printf("Ответ от Bitrix24 (start=%d): %s", start, string(responseData))

		var response struct {
			Result []map[string]string `json:"result"`
			Next   int                 `json:"next"`
		}

		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Println("Ошибка при разборе JSON:", err)
			return
		}

		contacts = append(contacts, response.Result...)
		log.Printf("Получено %d контактов, всего загружено: %d, следующий start: %v", len(response.Result), len(contacts), response.Next)

		if response.Next == 0 {
			break
		}

		start = response.Next
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("Всего загружено контактов с ДР за 5 дней до сегодня: %d", len(contacts))

	fileName := "birthDatesContactsLast5Days.json"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(contacts); err != nil {
		log.Println("Ошибка при записи в JSON-файл:", err)
		return
	}

	log.Println("Контакты с ДР за 5 дней до сегодня сохранены в", fileName)
}

// GetContactsListForScheduler Deprecated
func GetContactsListForScheduler() {
	webHookUrl := "https://klinikamedi.bitrix24.ru/rest/22/m0bcq858qrmfeoa5/crm.contact.list"

	contactsWithTodayBirthdays := make(map[string]string)
	start := 0

	// Получаем текущий день и месяц
	now := time.Now()
	currentDay := now.Day()
	currentMonth := int(now.Month())

	for {
		requestBody := map[string]interface{}{
			"FILTER": map[string]string{"!=BIRTHDATE": ""},
			"SELECT": []string{"ID", "BIRTHDATE"},
			"start":  start,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Println("Ошибка при маршалинге JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			return
		}
		defer resp.Body.Close()

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка при чтении ответа:", err)
			return
		}

		// Логируем ответ для диагностики
		log.Printf("Ответ от Bitrix24 (start=%d): %s", start, string(responseData))

		var response struct {
			Result []struct {
				ID        string `json:"ID"`
				BIRTHDATE string `json:"BIRTHDATE"`
			} `json:"result"`
			Next int `json:"next"`
		}

		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Println("Ошибка при разборе JSON:", err)
			return
		}

		// Проверяем, совпадает ли день рождения с текущей датой
		for _, result := range response.Result {
			parsedDate, err := time.Parse("2006-01-02T15:04:05-07:00", result.BIRTHDATE)
			if err != nil {
				log.Printf("Ошибка парсинга даты (%s): %v\n", result.BIRTHDATE, err)
				continue
			}

			if parsedDate.Day() == currentDay && int(parsedDate.Month()) == currentMonth {
				contactsWithTodayBirthdays[result.ID] = result.BIRTHDATE
			}
		}

		// Логируем текущее состояние
		log.Printf("Найдено %d контактов с сегодняшним днем рождения, всего обработано: %d, следующий start: %v", len(contactsWithTodayBirthdays), len(contactsWithTodayBirthdays), response.Next)

		// Проверяем, есть ли еще данные
		if response.Next == 0 {
			break
		}

		start = response.Next

		// Добавляем задержку, чтобы избежать лимитов API
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("Всего найдено контактов с днем рождения сегодня: %d", len(contactsWithTodayBirthdays))

	// Запись в JSON-файл (перезапись каждый день)
	fileName := "birthDates.json"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Форматированный JSON
	if err := encoder.Encode(contactsWithTodayBirthdays); err != nil {
		log.Println("Ошибка при записи в JSON-файл:", err)
		return
	}

	log.Println("Контакты с сегодняшним днем рождения сохранены в", fileName)
}

// GetContactsList Deprecated
func GetContactsList() {
	webHookUrl := "https://klinikamedi.bitrix24.ru/rest/22/m0bcq858qrmfeoa5/crm.contact.list"

	contactsWithBirthdays := make(map[string]string)
	start := 0

	for {
		requestBody := map[string]interface{}{
			"FILTER": map[string]string{"!=BIRTHDATE": "", "=UF_CRM_1739257086952": ""},
			"SELECT": []string{"ID", "BIRTHDATE"},
			"start":  start,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Println("Ошибка при маршалинге JSON:", err)
			return
		}

		req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			return
		}
		defer resp.Body.Close()

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка при чтении ответа:", err)
			return
		}

		// Логируем ответ для диагностики
		log.Printf("Ответ от Bitrix24 (start=%d): %s", start, string(responseData))

		var response struct {
			Result []struct {
				ID        string `json:"ID"`
				BIRTHDATE string `json:"BIRTHDATE"`
			} `json:"result"`
			Next int `json:"next"`
		}

		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Println("Ошибка при разборе JSON:", err)
			return
		}

		// Добавляем контакты в map
		for _, result := range response.Result {
			contactsWithBirthdays[result.ID] = result.BIRTHDATE
		}

		// Логируем текущее состояние
		log.Printf("Получено %d контактов, всего загружено: %d, следующий start: %v", len(response.Result), len(contactsWithBirthdays), response.Next)

		// Проверяем, есть ли еще данные
		if response.Next == 0 {
			break
		}

		start = response.Next

		// Добавляем задержку, чтобы избежать лимитов API
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("Всего загружено контактов: %d", len(contactsWithBirthdays))

	// Запись в JSON-файл
	fileName := "contacts.json"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Форматированный JSON
	if err := encoder.Encode(contactsWithBirthdays); err != nil {
		log.Println("Ошибка при записи в JSON-файл:", err)
		return
	}

	log.Println("Контакты сохранены в", fileName)
}

// Структура для тела запроса к Bitrix24
type UpdateContactRequest struct {
	ID     int `json:"id"`
	Fields struct {
		Day   int `json:"UF_CRM_1739257086952"` // День
		Month int `json:"UF_CRM_1739257138715"` // Месяц
		Year  int `json:"UF_CRM_1739257188787"` // Год
	} `json:"fields"`
}

// UpdateDateInContacts Deprecated
func UpdateDateInContacts() {
	fileName := "contacts.json"

	// Открываем JSON-файл
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Ошибка при чтении файла:", err)
		return
	}

	// Разбираем JSON в map
	contacts := make(map[string]string)
	if err := json.Unmarshal(data, &contacts); err != nil {
		log.Println("Ошибка при разборе JSON:", err)
		return
	}

	webHookUrl := "https://klinikamedi.bitrix24.ru/rest/22/m0bcq858qrmfeoa5/crm.contact.update"

	// Проходимся по всем контактам из JSON-файла
	for idStr, birthdate := range contacts {

		// Преобразуем ID из строки в число
		var id int
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			log.Printf("Ошибка преобразования ID (%s): %v\n", idStr, err)
			continue
		}

		// Парсим дату (формат: 1980-03-01T03:00:00+03:00)
		parsedDate, err := time.Parse("2006-01-02T15:04:05-07:00", birthdate)
		if err != nil {
			log.Printf("Ошибка парсинга даты (%s): %v\n", birthdate, err)
			continue
		}

		// Формируем JSON-запрос
		requestBody := UpdateContactRequest{ID: id}
		requestBody.Fields.Day = parsedDate.Day()
		requestBody.Fields.Month = int(parsedDate.Month())
		requestBody.Fields.Year = parsedDate.Year()

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Println("Ошибка при маршалинге JSON:", err)
			continue
		}

		// Отправляем запрос в Bitrix24
		req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

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

		// Логируем ответ
		log.Printf("Обновлен контакт ID=%d. Ответ Bitrix24: %s\n", id, strings.TrimSpace(string(responseData)))

		// Задержка, чтобы избежать лимитов API
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("Обновление контактов завершено.")
}
