package procurements

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	API_URL   = "https://torgi-crm.online/rest/88/672m3xnegwu7bosa/"
	HREF_FILE = "hrefs.json"
)

func AddBitrix(query map[string]interface{}, method string) (string, error) {
	url := API_URL + method
	body, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(responseData), nil
}

func GetResultBitrix(query map[string]interface{}, method string) ([]map[string]interface{}, error) {
	url := API_URL + method
	var results []map[string]interface{}
	next := 0
	total := 1

	for next <= total && total != 0 {
		query["start"] = next
		queryData, err := json.Marshal(query)
		if err != nil {
			return nil, err
		}

		resp, err := http.Post(url, "application/json", strings.NewReader(string(queryData)))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		if response["error_description"] != nil {
			return nil, fmt.Errorf("bitrix24 error: %s", response["error_description"])
		}

		if result, ok := response["result"].([]interface{}); ok {
			for _, item := range result {
				if m, ok := item.(map[string]interface{}); ok {
					results = append(results, m)
				}
			}
		}
		total = int(response["total"].(float64))
		next += 50
	}

	return results, nil
}

func GetInnData(inn, dateFrom string) (*goquery.Document, error) {
	time.Sleep(500 * time.Millisecond)
	url := fmt.Sprintf("https://zakupki.gov.ru/epz/order/extendedsearch/results.html?searchString=%s&morphology=on&search-filter=\u0414\u0410\u0422\u0415&pageNumber=1&sortDirection=false&recordsPerPage=_10&showLotsInfoHidden=false&sortBy=UPDATE_DATE&fz223=on&af=on&currencyIdGeneral=-1&publishDateFrom=%s", inn, dateFrom)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch URL: %s, status: %d", url, resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func RunProcurements() {
	start := time.Now()
	hrefsData, err := os.ReadFile(HREF_FILE)
	if err != nil {
		fmt.Printf("Error reading hrefs.json: %v\n", err)
		return
	}
	var allOrdersHrefs []string
	json.Unmarshal(hrefsData, &allOrdersHrefs)

	query := map[string]interface{}{
		"order":  map[string]string{"ID": "DESC"},
		"select": []string{"UF_INN", "ASSIGNED_BY_ID", "TITLE"},
	}

	companies, err := GetResultBitrix(query, "crm.company.list")
	if err != nil {
		fmt.Printf("Error fetching companies: %v\n", err)
		return
	}

	arUsers := map[string]string{}
	users, err := GetResultBitrix(query, "user.search")
	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		return
	}
	for _, user := range users {
		arUsers[user["ID"].(string)] = fmt.Sprintf("%s %s", user["NAME"], user["LAST_NAME"])
	}

	today := time.Now().Format("02.01.2006")
	todayNew := time.Now().AddDate(0, 0, -3).Format("02.01.2006")
	newCompaniesList := fmt.Sprintf("Отчёт от %s.\r\n", today)
	newCompaniesListArr := []string{}
	compCount := 0
	arrKey := 0

	for _, company := range companies {
		compInn, ok := company["UF_INN"].(string)
		if !ok || compInn == "" {
			continue
		}
		fmt.Printf("Processing INN: %s\n", compInn)
		html, err := GetInnData(compInn, todayNew)
		if err != nil {
			fmt.Printf("Error fetching INN data: %v\n", err)
			continue
		}

		newOrders := 0
		html.Find(".search-results__total").Each(func(i int, s *goquery.Selection) {
			newOrders, _ = strconv.Atoi(strings.TrimSpace(s.Text()))
		})

		if newOrders > 0 {
			if compCount%40 == 0 && compCount != 0 {
				arrKey++
			}
			manager := arUsers[fmt.Sprintf("%v", company["ASSIGNED_BY_ID"])]
			entry := fmt.Sprintf("%d) %s; ИНН - %s; Госзакупка - some_url - %s\r\n", compCount+1, company["TITLE"], compInn, manager)
			if len(newCompaniesListArr) <= arrKey {
				newCompaniesListArr = append(newCompaniesListArr, "")
			}
			newCompaniesListArr[arrKey] += entry
			allOrdersHrefs = append(allOrdersHrefs, fmt.Sprintf("some_url_for_%s", compInn))
			compCount++
		}
	}

	newCompaniesList += fmt.Sprintf("Новых закупок - %d\r\n", compCount)
	arrayMessage := map[string]interface{}{
		"DIALOG_ID": "chat6446",
		"MESSAGE":   newCompaniesList,
	}
	AddBitrix(arrayMessage, "im.message.add")

	for _, item := range newCompaniesListArr {
		arrayMessage := map[string]interface{}{
			"DIALOG_ID": "chat6446",
			"MESSAGE":   item,
		}
		AddBitrix(arrayMessage, "im.message.add")
		fileData, _ := json.MarshalIndent(allOrdersHrefs, "", "  ")
		os.WriteFile(HREF_FILE, fileData, 0644)
	}

	delta := time.Since(start)
	if delta.Seconds() > 30 {
		fmt.Printf("procurements: done with overload, execution time: %.2fs\n", delta.Seconds())
	} else {
		fmt.Printf("procurements: done, execution time: %.2fs\n", delta.Seconds())
	}
}
