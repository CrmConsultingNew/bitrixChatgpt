package torgi_new

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// BitrixResponse - структура ответа от Bitrix24
type BitrixResponse struct {
	Result []map[string]interface{} `json:"result"` // Правильная структура для JSON-массива
	Total  int                      `json:"total"`
	Next   int                      `json:"next"`
}

func GetCompaniesAndUsers() ([]map[string]interface{}, map[int]string) {
	// Фильтр для получения списка компаний
	companyFilters := map[string]interface{}{
		"order": map[string]string{"ID": "DESC"},
		"select": []string{
			"UF_INN",         // ИНН
			"ASSIGNED_BY_ID", // Ответственный
			"TITLE",          // Название
		},
	}

	// Фильтр для получения списка пользователей
	userFilters := map[string]interface{}{
		"order":  map[string]string{"ID": "ASC"},
		"select": []string{"ID", "NAME", "LAST_NAME"},
	}

	// Получаем список компаний
	arrayAllComp, err := GetResultBitrix(companyFilters, "crm.company.list")
	if err != nil {
		log.Fatalf("Error fetching company list: %v", err)
	}

	// Получаем список пользователей
	arUsersResult, err := GetResultBitrix(userFilters, "user.search")
	if err != nil {
		log.Fatalf("Error fetching user list: %v", err)
	}

	// Инициализируем карту для хранения пользователей
	arUsers := make(map[int]string)

	// Обрабатываем полученных пользователей
	for _, user := range arUsersResult {
		idStr, ok := user["ID"].(string) // ID приходит как строка
		if !ok {
			log.Println("Skipping user: invalid ID type")
			continue
		}

		// Конвертируем ID в int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Skipping user: failed to convert ID '%s' to int\n", idStr)
			continue
		}

		// Проверяем наличие имени и фамилии
		name, okName := user["NAME"].(string)
		lastName, okLastName := user["LAST_NAME"].(string)

		if okName && okLastName {

			//debug
			log.Printf("Добавлен пользователь: ID=%d, Имя=%s %s\n", id, name, lastName)

			arUsers[id] = fmt.Sprintf("%s %s", name, lastName)
		}
	}

	// Вывод результата для проверки
	fmt.Println("Полученные пользователи:")
	for id, fullName := range arUsers {
		fmt.Printf("ID: %d, Name: %s\n", id, fullName)
	}

	fmt.Println("Количество компаний:", len(arrayAllComp))

	//debug
	log.Printf("Передано в ProcessCompanies: %d компаний и %d пользователей\n", len(arrayAllComp), len(arUsers))

	// Возвращаем компании и пользователей
	return arrayAllComp, arUsers
}

// GetResultBitrix - Функция для выполнения запроса в Bitrix24
func GetResultBitrix(queryMain map[string]interface{}, method string) ([]map[string]interface{}, error) {
	fullURL := ApiURL + method
	next := 0
	total := 1
	var resultArray []map[string]interface{}

	for next <= total && total != 0 {
		queryMain["start"] = next

		// Преобразуем параметры в JSON
		queryBytes, err := json.Marshal(queryMain)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal query: %w", err)
		}

		// Выполняем HTTP-запрос
		resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(queryBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		// Читаем ответ
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		// Разбираем JSON
		var response BitrixResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}

		// Проверяем, не пустой ли ответ
		if len(response.Result) == 0 {
			break // Если пусто - выходим
		}

		// Добавляем полученные данные в общий результат
		resultArray = append(resultArray, response.Result...)
		total = response.Total
		next += 50
	}

	return resultArray, nil
}

func StartFuncGetResultBitrixAtFirst() {
	queryMain := map[string]interface{}{
		"order": map[string]string{"ID": "DESC"},
		"select": []string{
			"UF_INN",
			"ASSIGNED_BY_ID",
			"TITLE",
		},
	}

	results, err := GetResultBitrix(queryMain, "crm.company.list")
	if err != nil {
		fmt.Printf("Error fetching data from Bitrix24: %v\n", err)
		return
	}

	fmt.Printf("Results: %+v\n", results)
}
