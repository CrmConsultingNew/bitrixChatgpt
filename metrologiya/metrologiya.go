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

type ReportField struct {
	Key         string
	Description string
}

// Упорядоченные поля для отчетов
var reportFields = []ReportField{
	{Key: "UF_CRM_1717764898", Description: "Поверка газ. счетчика"}, // Переделать на Поверка газовых счетчиков
	{Key: "UF_CRM_1720091191105", Description: "Диагностика газ. котла"},
	{Key: "UF_CRM_1717765382", Description: "Диагностика газовой плиты"},
	{Key: "UF_CRM_1717765494", Description: "Поверка водяного счетчика"},
	{Key: "UF_CRM_1729173539594", Description: "Диагностика промышленных котлов"},
	{Key: "UF_CRM_1717765416", Description: "Продажа шланга"},
	{Key: "UF_CRM_1717765660", Description: "Продажа газового счетчика"}, // Переделать на Продажа газовых счетчиков
	{Key: "UF_CRM_1718280444", Description: "Продажа водяного счетчика"},
	{Key: "UF_CRM_1725598657143", Description: "QR-код"},
	{Key: "UF_CRM_1718021528233", Description: "Безналичный перевод"},
	{Key: "UF_CRM_1716890120764", Description: "Сумма компании"},
	{Key: "UF_CRM_1716890078180", Description: "Процент техника поверителя"},
	{Key: "UF_CRM_1716890031608", Description: "ГСМ"},
	{Key: "UF_CRM_1717765798", Description: "Квартира за сутки"},
	{Key: "UF_CRM_1717765861", Description: "Прочие расходы"},
	{Key: "UF_CRM_1718022170", Description: "Скидка"}, // Новое поле
}

// IDs ответственных сотрудников
var AssignedByIDs []string

// Форматируем дату в ISO 8601
func formatDateToISO(date string) string {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Fatalf("Ошибка парсинга даты: %v", err)
	}
	return parsedDate.Format("2006-01-02T00:00:00+03:00")
}

// Выборка по Стадиям
func determineStageID(fieldKey string) []string {
	switch fieldKey {
	case "UF_CRM_1725598657143", "UF_CRM_1718021528233", "UF_CRM_1716890120764",
		"UF_CRM_1716890078180", "UF_CRM_1716890031608", "UF_CRM_1717765798", "UF_CRM_1717765861",
		"UF_CRM_1718022170": // Для поля "Скидка"
		return []string{"WON", "EXECUTING"}
	case "UF_CRM_1717764898": // Поверка газгового счетчика
		return []string{"EXECUTING"}
	default:
		return []string{"WON"}
	}
}

// Общий обработчик для отправки отчетов
func generateReport(startDate, endDate, reportTitle string, chatID string, isDaily bool) {

	reportDates := ""
	if endDate != "" {
		reportDates = fmt.Sprintf(" от %s до %s", startDate, endDate)
	} else {
		reportDates = fmt.Sprintf(" за %s", startDate)
	}

	resultText := reportTitle + reportDates + "\n"
	for _, field := range reportFields {
		var total int
		var sumOfField float64
		stageIDs := determineStageID(field.Key)
		if isDaily {
			total, sumOfField = processField(field.Key, field.Description, stageIDs, startDate)
		} else {
			total, sumOfField = processFieldRange(field.Key, field.Description, stageIDs, startDate, endDate)
		}
		resultText += fmt.Sprintf("%s: Всего: %d Сумма: %.2f.\n", field.Description, total, sumOfField)
	}

	err := sendMessageToBitrixChat(chatID, resultText)
	if err != nil {
		log.Printf("Ошибка отправки сообщения в чат Bitrix24: %v", err)
	}
}

func processField(field, description string, stageIDs []string, date string) (int, float64) {
	log.Printf("Processing field: %s, Stage IDs: %v, Date: %s", description, stageIDs, date)

	filter := map[string]interface{}{
		"=STAGE_ID":               stageIDs,
		"=CATEGORY_ID":            0,
		"ASSIGNED_BY_ID":          AssignedByIDs,
		"=UF_CRM_1650442073068":   date,
		fmt.Sprintf(">%s", field): 0,
	}

	total, sum := fetchData(field, description, filter)

	log.Printf("Field: %s, Total: %d, Sum: %.2f", description, total, sum)

	return total, sum
}

func processFieldRange(field, description string, stageIDs []string, startDate, endDate string) (int, float64) {
	log.Printf("Processing field range: %s, Stage IDs: %v, StartDate: %s, EndDate: %s", description, stageIDs, startDate, endDate)

	filter := map[string]interface{}{
		"=STAGE_ID":               stageIDs,
		"=CATEGORY_ID":            0,
		"ASSIGNED_BY_ID":          AssignedByIDs,
		">=UF_CRM_1650442073068":  startDate,
		"<=UF_CRM_1650442073068":  endDate,
		fmt.Sprintf(">%s", field): 0,
	}

	total, sum := fetchData(field, description, filter)

	log.Printf("Field: %s, Total: %d, Sum: %.2f", description, total, sum)

	return total, sum
}

// fetchData отправляет запрос и возвращает результат
func fetchData(field, description string, filter map[string]interface{}) (int, float64) {
	requestBody := map[string]interface{}{
		"filter": filter,
		"select": []string{"ID", "TITLE", "STAGE_ID", "CATEGORY_ID", field},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Ошибка сериализации запроса для %s: %v", description, err)
	}

	log.Printf("Отправляемый запрос для %s: %s", description, string(jsonData))

	req, err := http.NewRequest("POST", webHookUrl+"crm.deal.list", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Ошибка создания запроса для %s: %v", description, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса для %s: %v", description, err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения ответа для %s: %v", description, err)
	}

	log.Printf("Ответ сервера для %s: %s", description, string(responseData))

	var response struct {
		Result []map[string]interface{} `json:"result"`
		Total  int                      `json:"total"`
	}

	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Fatalf("Ошибка декодирования JSON-ответа для %s: %v", description, err)
	}

	var dealsIds []string
	sum := 0.0
	for _, deal := range response.Result {
		if id, ok := deal["ID"].(string); ok {
			dealsIds = append(dealsIds, id)
		}
		if fieldValue, ok := deal[field].(string); ok {
			var value float64
			fmt.Sscanf(fieldValue, "%f", &value)
			sum += value
		} else if fieldValueFloat, ok := deal[field].(float64); ok {
			sum += fieldValueFloat
		}
	}

	log.Println("dealsIds: ", dealsIds)

	bizprocMap, err := FetchBizprocInstances()
	if err != nil {
		log.Fatalf("Ошибка получения данных от bizproc.workflow.instances: %v", err)
	}

	var bizprocsIds []string
	for _, dealID := range dealsIds {
		for bizprocID, docID := range bizprocMap {
			if dealID == docID {
				bizprocsIds = append(bizprocsIds, bizprocID)
			}
		}
	}

	log.Printf("Сформированный массив bizprocsIds: %v", bizprocsIds)

	for _, bizprocID := range bizprocsIds {
		if err := TerminateBizproc(bizprocID); err != nil {
			log.Printf("Ошибка завершения bizproc.workflow: %v", err)
		}
	}

	// Новый блок для запуска workflow через bizproc.workflow.start
	for _, dealID := range dealsIds {
		if err := StartBizprocWorkflow("295", dealID); err != nil {
			log.Printf("Ошибка запуска bizproc.workflow.start для сделки %s: %v", dealID, err)
		}
	}

	return response.Total, sum
}

// StartBizprocWorkflow запускает workflow для заданного TEMPLATE_ID и DOCUMENT_ID
func StartBizprocWorkflow(templateID string, dealID string) error {
	requestBody := map[string]interface{}{
		"TEMPLATE_ID": templateID,
		"DOCUMENT_ID": []string{"crm", "CCrmDocumentDeal", "DEAL_" + dealID},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("Ошибка сериализации запроса для bizproc.workflow.start: %v", err)
	}

	req, err := http.NewRequest("POST", webHookUrl+"bizproc.workflow.start", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Ошибка создания запроса для bizproc.workflow.start: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Ошибка выполнения запроса для bizproc.workflow.start: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Ошибка чтения ответа для bizproc.workflow.start: %v", err)
	}

	log.Printf("Ответ сервера для bizproc.workflow.start (dealID=%s): %s", dealID, string(responseData))
	return nil
}

const webHookUrl = "https://metrologiya.bitrix24.ru/rest/93/tbpyex30x7a5d4gs/"

// FetchBizprocInstances retrieves all data from bizproc.workflow.instances, handling pagination.
func FetchBizprocInstances() (map[string]string, error) {
	bizprocMap := make(map[string]string)
	start := 0

	for {
		// Формирование тела запроса
		requestBody := map[string]interface{}{
			"select": []string{"ID", "DOCUMENT_ID"},
			"start":  start,
		}

		// Сериализация тела запроса
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сериализации запроса для bizproc.workflow.instances: %v", err)
		}

		// Создание HTTP-запроса
		req, err := http.NewRequest("POST", webHookUrl+"bizproc.workflow.instances", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("Ошибка создания запроса для bizproc.workflow.instances: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Выполнение запроса
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Ошибка выполнения запроса для bizproc.workflow.instances: %v", err)
		}
		defer resp.Body.Close()

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Ошибка чтения ответа для bizproc.workflow.instances: %v", err)
		}

		// Логируем ответ для отладки
		log.Printf("Ответ сервера (start=%d): %s", start, string(responseData))

		// Декодирование ответа
		var response struct {
			Result []map[string]interface{} `json:"result"`
			Next   int                      `json:"next"`
			Total  int                      `json:"total"`
		}

		if err := json.Unmarshal(responseData, &response); err != nil {
			return nil, fmt.Errorf("Ошибка декодирования JSON-ответа для bizproc.workflow.instances: %v", err)
		}

		// Обработка полученных данных
		for _, instance := range response.Result {
			if id, ok := instance["ID"].(string); ok {
				if documentID, ok := instance["DOCUMENT_ID"].(string); ok {
					if len(documentID) > 5 && documentID[:5] == "DEAL_" {
						bizprocMap[id] = documentID[5:] // Обрезаем префикс "DEAL_"
					}
				}
			}
		}

		// Проверяем, есть ли следующая страница
		if response.Next == 0 {
			break
		}

		start = response.Next // Обновляем значение для следующего запроса
	}

	return bizprocMap, nil
}

// TerminateBizproc terminates a bizproc.workflow by ID
func TerminateBizproc(bizprocID string) error {
	requestBody := map[string]interface{}{
		"ID": bizprocID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("Ошибка сериализации запроса для bizproc.workflow.terminate: %v", err)
	}

	req, err := http.NewRequest("POST", webHookUrl+"bizproc.workflow.terminate", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Ошибка создания запроса для bizproc.workflow.terminate: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Ошибка выполнения запроса для bizproc.workflow.terminate: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("bizproc.workflow.terminate выполнен для ID: %s", bizprocID)
	return nil
}
