package reports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func curlPOST(endpoint string, query url.Values) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.PostForm(endpoint, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func callB24Method(method string, params map[string]interface{}) (map[string]interface{}, error) {
	if B24SLEEP > 0 {
		time.Sleep(B24SLEEP)
	}

	webhook := fmt.Sprintf("https://%s/rest/%s/%s/", B24_HOST, B24_USER_ID, B24_WEBHOOK)
	endpoint := webhook + method + ".json"

	query := url.Values{}
	for key, value := range params {
		query.Set(key, fmt.Sprintf("%v", value))
	}

	response, err := curlPOST(endpoint, query)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	logMessage("B24 REQUEST METHOD", method)
	logMessage("B24 REQUEST PARAMS", params)
	logMessage("B24 RESPONSE", result)

	return result, nil
}

func callB24Batch(cmd []map[string]interface{}) (map[string]interface{}, error) {
	batch := map[string]interface{}{
		"halt": "0",
		"cmd":  make(map[string]string),
	}

	for i, item := range cmd {
		params := url.Values{}
		for key, value := range item["params"].(map[string]interface{}) {
			params.Set(key, fmt.Sprintf("%v", value))
		}
		batch["cmd"].(map[string]string)[strconv.Itoa(i)] = fmt.Sprintf("%s?%s", item["method"], params.Encode())
	}

	return callB24Method("batch", batch)
}

func getFullList(method string, params map[string]interface{}) ([]interface{}, error) {
	result, err := callB24Method(method, params)
	if err != nil {
		return nil, err
	}

	list := result["result"].([]interface{})
	cmd := []map[string]interface{}{}

	for i := 50; i < int(result["total"].(float64)); i += 50 {
		params["start"] = i
		cmd = append(cmd, map[string]interface{}{
			"method": method,
			"params": params,
		})
	}

	if len(cmd) > 0 {
		batchResult, err := callB24Batch(cmd)
		if err != nil {
			return nil, err
		}
		for _, page := range batchResult["result"].(map[string]interface{})["result"].([]interface{}) {
			list = append(list, page.([]interface{})...)
		}
	}

	return list, nil
}

func getCallList(filter map[string]interface{}) ([]interface{}, error) {
	return getFullList("voximplant.statistic.get", map[string]interface{}{
		"FILTER": filter,
		"SORT":   "CALL_START_DATE",
		"ORDER":  "ASC",
	})
}
func GetCallList(filter map[string]interface{}) ([]map[string]interface{}, error) {
	data, err := getCallList(filter)
	if err != nil {
		return nil, err
	}
	// Преобразование интерфейса в карту
	var result []map[string]interface{}
	for _, item := range data {
		result = append(result, item.(map[string]interface{}))
	}
	return result, nil
}

func getLeadList(filter, order map[string]interface{}) ([]interface{}, error) {
	return getFullList("crm.lead.list", map[string]interface{}{
		"FILTER": filter,
		"ORDER":  order,
	})
}

func getDealList(filter, order map[string]interface{}) ([]interface{}, error) {
	return getFullList("crm.deal.list", map[string]interface{}{
		"FILTER": filter,
		"ORDER":  order,
	})
}

func getTaskList(filter, order map[string]interface{}) ([]interface{}, error) {
	return getFullList("task.item.list", map[string]interface{}{
		"ORDER":  order,
		"FILTER": filter,
		"PARAMS": map[string]string{"LOAD_TAGS": "Y"},
		"SELECT": []string{},
	})
}

func getUserList(filter map[string]interface{}) ([]interface{}, error) {
	return getFullList("user.get", map[string]interface{}{
		"FILTER": filter,
	})
}

func logMessage(title string, data interface{}) {
	fmt.Printf("%s: %v\n", title, data)
}
