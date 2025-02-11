package torgi_new

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// ProcessNewCompanies Основной процесс обработки новых закупок
func ProcessNewCompanies(today string, compCount int, allOrdersHrefs []string, newCompaniesListArr map[int]string) {
	// Формируем общий список закупок
	newCompaniesList := fmt.Sprintf("Отчёт от %s.\r\nНовых закупок - %d\r\n", today, compCount)

	// Создаем сообщение для отправки в Bitrix
	arrayMessage := map[string]string{
		"DIALOG_ID": "chat6446",
		"MESSAGE":   newCompaniesList,
	}

	// Сохраняем обновленный список ссылок в файл
	err := SaveHrefs(HrefFile, allOrdersHrefs)
	if err != nil {
		fmt.Printf("Ошибка сохранения hrefs.json: %s\n", err)
		return
	}

	// Отправляем общее сообщение в Bitrix
	response, err := AddBitrix(arrayMessage, "im.message.add")
	if err != nil {
		fmt.Printf("Ошибка отправки сообщения в Bitrix: %s\n", err)
		return
	}
	fmt.Printf("Ответ от Bitrix: %s\n", response)

	// Обрабатываем каждую группу закупок и отправляем в Bitrix
	for _, item := range newCompaniesListArr {
		newCompaniesList = item
		arrayMessage := map[string]string{
			"DIALOG_ID": "chat6446",
			"MESSAGE":   newCompaniesList,
		}

		// Сохраняем обновленный список ссылок
		err := SaveHrefs(HrefFile, allOrdersHrefs)
		if err != nil {
			fmt.Printf("Ошибка сохранения hrefs.json: %s\n", err)
			continue
		}

		// Отправляем сообщение в Bitrix
		response, err := AddBitrix(arrayMessage, "im.message.add")
		if err != nil {
			fmt.Printf("Ошибка отправки сообщения в Bitrix: %s\n", err)
			continue
		}
		fmt.Printf("Ответ от Bitrix: %s\n", response)
	}
}

// SaveHrefs Функция для сохранения hrefs.json
func SaveHrefs(filename string, hrefs []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(hrefs); err != nil {
		return fmt.Errorf("не удалось сохранить данные: %w", err)
	}
	return nil
}

// Общий метод AddBitrix
func AddBitrix(queryMain map[string]string, method string) (string, error) {
	// Формируем URL с методом
	fullURL := MyApiUrl + method

	// Преобразуем параметры в URL-encoded строку
	data := url.Values{}
	for key, value := range queryMain {
		data.Set(key, value)
	}

	// Создаём HTTP-запрос
	resp, err := http.Post(fullURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ от сервера
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}
