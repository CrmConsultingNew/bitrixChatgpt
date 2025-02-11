package metrologiya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchUserIDsByDepartment() ([]string, error) {
	webHookUrl := "https://metrologiya.bitrix24.ru/rest/93/tbpyex30x7a5d4gs/user.search"
	requestBody := map[string]interface{}{
		"UF_DEPARTMENT_NAME": "поверители",
	}
	// ранее было так "WORK_POSITION": "поверитель"
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %v", err)
	}

	req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания HTTP-запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения HTTP-запроса: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	log.Printf("Ответ от user.search: %s", string(responseData))

	var response struct {
		Result []struct {
			ID string `json:"ID"`
		} `json:"result"`
	}

	if err := json.Unmarshal(responseData, &response); err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON-ответа: %v", err)
	}

	for _, user := range response.Result {
		AssignedByIDs = append(AssignedByIDs, user.ID)
		AssignedByIDs = append(AssignedByIDs, "9347")
	}

	log.Printf("AssignedByIDs: %v", AssignedByIDs)

	return AssignedByIDs, nil
}
