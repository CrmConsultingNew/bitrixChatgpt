package mail

import (
	"encoding/json"
	"log"
	"net/http"
)

type MailRequest struct {
	Email          string `json:"email"`
	Body           string `json:"body"`
	AttachmentPath string `json:"attachmentPath"` // Путь к файлу вложения
}

func SendMailHandler(w http.ResponseWriter, r *http.Request) {
	// Декодирование JSON-запроса для получения email и body
	var mailReq MailRequest
	if err := json.NewDecoder(r.Body).Decode(&mailReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Устанавливаем фиксированный путь к файлу
	mailReq.AttachmentPath = "/root/bitrixChatgpt/tables.docx"

	// Конфигурация SMTP
	smtpConfig := SMTPConfig{
		Host:     "smtp.beget.com",
		Port:     "465",
		Username: "crm@crmconsulting-api.ru", // Полный адрес электронной почты
		Password: "AA1379657aa!",
		From:     "crm@crmconsulting-api.ru", // Полный адрес электронной почты как отправитель
	}

	// Добавляем текст в тело письма
	mailReq.Body = "Отправляем акт-сверки от компании CRM-Consulting. Акт сверки во вложении. Чтобы файл корректно отображался - скачайте файл."

	// Отправка письма с вложением
	err := SendEmail(smtpConfig, mailReq.Email, mailReq.Body, mailReq.AttachmentPath)
	if err != nil {
		log.Println("SendEmail failed: ", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
