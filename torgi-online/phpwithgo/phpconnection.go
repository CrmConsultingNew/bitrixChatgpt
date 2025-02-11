package phpwithgo

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type ResponseData struct {
	Success     bool   `json:"success"`
	TotalOrders int    `json:"total_orders"`
	CompHref    string `json:"comp_href"`
	DateStart   string `json:"date_start"`
	DateUpdate  string `json:"date_update"`
	DateFinish  string `json:"date_finish"`
}

func parseHTML(inn string, dateFrom string) (*ResponseData, error) {
	url := fmt.Sprintf("https://zakupki.gov.ru/epz/order/extendedsearch/results.html?searchString=%s&morphology=on&search-filter=Дате+размещения&publishDateFrom=%s", inn, dateFrom)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Парсим .search-results__total
	totalOrders := 0
	doc.Find(".search-results__total").Each(func(i int, s *goquery.Selection) {
		fmt.Sscanf(s.Text(), "%d", &totalOrders)
	})

	// Парсим .registry-entry__header-mid__number
	compHref := ""
	doc.Find(".registry-entry__header-mid__number a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			compHref = "https://zakupki.gov.ru" + href
		}
	})

	// Парсим .data-block__value
	dates := []string{}
	doc.Find(".data-block__value").Each(func(i int, s *goquery.Selection) {
		dates = append(dates, s.Text())
	})

	if len(dates) < 3 {
		return nil, fmt.Errorf("failed to extract all dates")
	}

	return &ResponseData{
		Success:     true,
		TotalOrders: totalOrders,
		CompHref:    compHref,
		DateStart:   dates[0],
		DateUpdate:  dates[1],
		DateFinish:  dates[2],
	}, nil
}

func TorgiOnlineHandler(w http.ResponseWriter, r *http.Request) {
	inn := r.URL.Query().Get("inn")
	dateFrom := r.URL.Query().Get("date_from")
	if inn == "" || dateFrom == "" {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	data, err := parseHTML(inn, dateFrom)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
