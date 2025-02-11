package metrologiya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Bizproc struct {
	ID       string `json:"ID"`
	Modified string `json:"MODIFIED"`
}

type Response struct {
	Result []Bizproc `json:"result"`
	Next   int       `json:"next"`
	Total  int       `json:"total"`
}

func GetAllBizProcData() {
	webHookUrl := "https://metrologiya.bitrix24.ru/rest/93/why1zgds7x86kf9a/bizproc.workflow.instances"
	allBizProcs := []Bizproc{}
	start := 0

	// Определяем диапазон дат: от 3 месяцев назад до текущей даты
	now := time.Now()
	startDate := now.AddDate(0, -3, 0).Format("2006-01-02T15:04:05Z07:00")
	endDate := now.AddDate(0, 0, -1).Format("2006-01-02T15:04:05Z07:00")

	for {
		// Формируем запрос с параметром пагинации и фильтром по дате
		requestBody := map[string]interface{}{
			"start": start,
			"FILTER": map[string]interface{}{
				">=MODIFIED": startDate,
				"<=MODIFIED": endDate,
			},
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatalf("Ошибка сериализации запроса: %v", err)
		}

		req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Ошибка создания запроса: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("Ошибка выполнения запроса: %v", err)
		}
		defer resp.Body.Close()

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Ошибка чтения ответа: %v", err)
		}

		var response Response
		if err := json.Unmarshal(responseData, &response); err != nil {
			log.Fatalf("Ошибка декодирования JSON-ответа: %v", err)
		}

		// Добавляем результаты текущего запроса в общий список
		allBizProcs = append(allBizProcs, response.Result...)

		// Если нет больше данных для получения, выходим из цикла
		if response.Next == 0 {
			break
		}

		// Устанавливаем start для следующего запроса
		start += 50
	}

	// Выводим общее количество записей
	fmt.Printf("Общее количество бизнес-процессов за последние 3 месяца: %d\n", len(allBizProcs))

	// Перебираем бизнес-процессы и вызываем StopAllBizProcs для каждого
	for _, bizproc := range allBizProcs {
		time.Sleep(time.Second * 1) // Добавляем паузу между запросами
		StopAllBizProcs(bizproc.ID)
	}
}

func StopAllBizProcs(ID string) {
	webHookUrl := "https://metrologiya.bitrix24.ru/rest/93/why1zgds7x86kf9a/bizproc.workflow.terminate"

	// Формируем запрос с параметром для завершения бизнес-процесса
	requestBody := map[string]interface{}{
		"ID": ID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Ошибка сериализации запроса: %v", err)
	}

	req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения ответа: %v", err)
	}

	// Выводим ответ для отладки
	log.Println("Ответ от Bitrix24:", string(responseData))
}
