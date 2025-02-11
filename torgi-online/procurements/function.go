package procurements

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type OrderPlan struct {
	FileName        string
	FileNameCompany string
	Date            string
	ArResult        map[string]interface{}
}

func NewOrderPlan(fileName string) *OrderPlan {
	return &OrderPlan{
		FileName:        fileName,
		FileNameCompany: "company.json",
		Date:            time.Now().Format("02.01.2006"),
		ArResult:        make(map[string]interface{}),
	}
}

func (op *OrderPlan) GetCurl(url string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	randomChromeVersion := fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.%d Safari/537.36",
		rand.Intn(3)+111, rand.Intn(50))
	req.Header.Set("User-Agent", randomChromeVersion)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (op *OrderPlan) ParseDOM(data string) (int, string, string, error) {
	totalRegex := regexp.MustCompile(`<div class="search-results__total">(\d+)</div>`)
	totalMatch := totalRegex.FindStringSubmatch(data)

	if len(totalMatch) < 2 {
		return 0, "", "", fmt.Errorf("failed to parse .search-results__total")
	}

	total, err := strconv.Atoi(totalMatch[1])
	if err != nil {
		return 0, "", "", err
	}

	// Example: Extend for `.data-block__value` or `.registry-entry__header-mid__number`
	return total, "start_date", "end_date", nil
}

func (op *OrderPlan) GetJson(fileName string) (map[string]interface{}, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (op *OrderPlan) SetJson(data map[string]interface{}, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}
