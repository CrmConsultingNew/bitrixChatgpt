package mail

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/smtp"
	"path/filepath"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// LOGINAuth creates an smtp.Auth that implements the LOGIN authentication mechanism.
func LOGINAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unexpected server challenge: %s", fromServer)
		}
	}
	return nil, nil
}

// TestSMTPConnection tests the connection to the SMTP server with the provided configuration.
func TestSMTPConnection(smtpConfig SMTPConfig) error {
	addr := fmt.Sprintf("%s:%s", smtpConfig.Host, smtpConfig.Port)

	log.Println("Connection to SMTP was started")

	// Create a tls.Config with the ServerName set to the SMTP host
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpConfig.Host,
	}

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Create new SMTP client
	client, err := smtp.NewClient(conn, smtpConfig.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Authenticate using LOGIN
	auth := LOGINAuth(smtpConfig.Username, smtpConfig.Password)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	log.Println("SMTP connection successful")
	return nil
}

// SendEmail отправляет письмо с вложением на указанный email.
func SendEmail(smtpConfig SMTPConfig, toEmail, body, attachmentPath string) error {
	// Формируем заголовки сообщения
	subject := "Important Message"
	header := make(map[string]string)
	header["From"] = smtpConfig.From
	header["To"] = toEmail
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `multipart/mixed; boundary="simpleboundary"`

	// Формируем тело сообщения с вложением
	message := strings.Builder{}
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n--simpleboundary\r\n")
	message.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	message.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	message.WriteString(body)
	message.WriteString("\r\n--simpleboundary\r\n")

	// Если указан путь к файлу вложения, добавляем его в письмо
	if attachmentPath != "" {
		fileData, err := ioutil.ReadFile(attachmentPath)
		if err != nil {
			return fmt.Errorf("failed to read attachment file: %v", err)
		}

		filename := filepath.Base(attachmentPath)
		message.WriteString(fmt.Sprintf("Content-Type: application/octet-stream; name=\"%s\"\r\n", mime.QEncoding.Encode("UTF-8", filename)))
		message.WriteString("Content-Transfer-Encoding: base64\r\n")
		message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", mime.QEncoding.Encode("UTF-8", filename)))

		encoded := base64.StdEncoding.EncodeToString(fileData)
		message.WriteString(encoded)
		message.WriteString("\r\n--simpleboundary--\r\n")
	}

	// Адрес сервера SMTP
	addr := fmt.Sprintf("%s:%s", smtpConfig.Host, smtpConfig.Port)

	// Настройка TLS
	tlsConfig := &tls.Config{
		ServerName: smtpConfig.Host,
	}

	// Подключение к SMTP серверу
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Создание SMTP клиента
	client, err := smtp.NewClient(conn, smtpConfig.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Аутентификация
	auth := LOGINAuth(smtpConfig.Username, smtpConfig.Password)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Установка отправителя
	if err := client.Mail(smtpConfig.From); err != nil {
		return fmt.Errorf("failed to set sender address: %v", err)
	}

	// Установка получателя
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("failed to set recipient address: %v", err)
	}

	// Отправка сообщения
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send email body: %v", err)
	}

	_, err = wc.Write([]byte(message.String()))
	if err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close email writer: %v", err)
	}

	return nil
}
