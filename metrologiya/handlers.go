package metrologiya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Отправка сообщения в чат Bitrix24
func sendMessageToBitrixChat(dialogID, message string) error {
	webHookUrl := "https://metrologiya.bitrix24.ru/rest/93/tbpyex30x7a5d4gs/im.message.add"
	requestBody := map[string]interface{}{"DIALOG_ID": dialogID, "MESSAGE": message}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("ошибка сериализации тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания HTTP-запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения HTTP-запроса: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// Обработчик для вызова ежедневного отчета
func HandleDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	generateReport(yesterday, "", "Ежедневный отчет", "chat39259", true)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ежедневный отчет успешно отправлен"))
}

// Обработчик для вызова еженедельного отчета
func HandleWeeklyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	endDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	generateReport(startDate, endDate, "Еженедельный отчет", "chat39259", false)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Еженедельный отчет успешно отправлен"))
}

// Обработчик для вызова ежемесячного отчета
func HandleMonthlyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	firstDayOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDayOfPreviousMonth := firstDayOfCurrentMonth.AddDate(0, 0, -1)
	firstDayOfPreviousMonth := time.Date(lastDayOfPreviousMonth.Year(), lastDayOfPreviousMonth.Month(), 1, 0, 0, 0, 0, now.Location())

	startDate := firstDayOfPreviousMonth.Format("2006-01-02")
	endDate := lastDayOfPreviousMonth.Format("2006-01-02")

	generateReport(startDate, endDate, "Ежемесячный отчет", "chat39259", false)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ежемесячный отчет успешно отправлен"))
}
