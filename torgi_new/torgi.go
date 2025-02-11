package torgi_new

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//MyApiUrl = "https://b24-iesb30.bitrix24.ru/rest/1/ytdf4fz89gmf7wsp/"

const (
	ApiURL   = "https://torgi-crm.online/rest/88/672m3xnegwu7bosa/"
	MyApiUrl = "https://torgi-crm.online/rest/88/672m3xnegwu7bosa/"
	HrefFile = "torgi_new/hrefs.json"
)

// GetInn отправляет запрос к указанному URL с параметрами `compInn` и `todayNew`
// и возвращает HTML-страницу в виде строки. То есть ищет на странице ответы по ИНН, должна использоваться после поиска ИНН Компаний в Битрикс
func GetInn(compInn string, todayNew string) (string, error) {
	// Задержка перед отправкой запроса (аналог `usleep(500000)` в PHP)
	time.Sleep(500 * time.Millisecond)

	// Формирование URL с параметрами
	url := fmt.Sprintf("https://zakupki.gov.ru/epz/order/extendedsearch/results.html?searchString=%s&morphology=on&search-filter=%%D0%%94%%D0%%B0%%D1%%82%%D0%%B5+%%D1%%80%%D0%%B0%%D0%%B7%%D0%%BC%%D0%%B5%%D1%%89%%D0%%B5%%D0%%BD%%D0%%B8%%D1%%8F&pageNumber=1&sortDirection=false&recordsPerPage=_10&showLotsInfoHidden=false&sortBy=UPDATE_DATE&fz223=on&af=on&currencyIdGeneral=-1&publishDateFrom=%s", compInn, todayNew)

	// Создание HTTP-запроса
	client := &http.Client{
		Timeout: 15 * time.Second, // Устанавливаем тайм-аут в 15 секунд
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем пользовательский агент
	userAgent := fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.%d Safari/537.36",
		113+RandInt(0, 1), RandInt(0, 50)) // Генерация номера версии, как в PHP
	req.Header.Set("User-Agent", userAgent)

	// Отправка запроса
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем HTTP-статус
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	// Читаем содержимое страницы
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// RandInt генерирует случайное целое число в указанном диапазоне.
func RandInt(min, max int) int {
	return min + int(time.Now().UnixNano()%int64(max-min+1))
}

// LoadHrefs загружает JSON-данные из файла hrefs.json и возвращает слайс строк.
func LoadHrefs(filename string) []string {
	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filename, err)
		return []string{}
	}
	defer file.Close() // Закрываем файл в конце

	// Читаем JSON-данные
	var hrefs []string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&hrefs); err != nil {
		log.Printf("Failed to decode JSON from file %s: %v", filename, err)
		return []string{}
	}

	return hrefs
}

func StartFuncLoadHrefs() {
	allOrdersHrefs := LoadHrefs(HrefFile)
	if len(allOrdersHrefs) == 0 {
		log.Println("No hrefs loaded or file is empty.")
	} else {
		log.Printf("Loaded %d hrefs from file.\n", len(allOrdersHrefs))
	}
}
