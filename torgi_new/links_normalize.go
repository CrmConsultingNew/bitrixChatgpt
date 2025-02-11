package torgi_new

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func StartNormalizeLinks() {
	inputFile := "torgi_new/hrefs_old.json"
	outputFile := "torgi_new/hrefs.json"

	// Читаем ссылки из файла
	urls, err := ReadURLsFromFile(inputFile)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	// Нормализуем ссылки
	normalizedURLs := NormalizeURLs(urls)

	// Записываем результат в новый файл
	err = WriteURLsToFile(outputFile, normalizedURLs)
	if err != nil {
		log.Fatalf("Ошибка записи файла: %v", err)
	}

	log.Println("✅ Ссылки успешно нормализованы и сохранены в hrefs.json")
}

// NormalizeURLs обрабатывает и очищает URL-адреса
func NormalizeURLs(urls []string) []string {
	var normalizedURLs []string

	for _, url := range urls {
		// Убираем лишние слэши, но сохраняем корректный формат https://
		cleanURL := strings.Replace(url, "https:/", "https://", 1) // Фиксим https://
		cleanURL = strings.Replace(cleanURL, "//", "/", -1)        // Убираем двойные слэши

		normalizedURLs = append(normalizedURLs, cleanURL)
	}

	return normalizedURLs
}

// ReadURLsFromFile читает список URL из JSON-файла
func ReadURLsFromFile(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %s: %v", filename, err)
	}

	var urls []string
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return urls, nil
}

// WriteURLsToFile записывает нормализованные URL в новый JSON-файл
func WriteURLsToFile(filename string, urls []string) error {
	data, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка кодирования JSON: %v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла %s: %v", filename, err)
	}

	return nil
}
