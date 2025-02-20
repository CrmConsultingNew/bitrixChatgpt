package reportsDDS

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	bitrixAPI        = "https://crmconsalting.bitrix24.ru/rest/1536/3pllzln0z4a69b23/lists.element.get"
	bitrixDealAPI    = "https://crmconsalting.bitrix24.ru/rest/1536/3pllzln0z4a69b23/crm.deal.get"
	bitrixCompanyAPI = "https://crmconsalting.bitrix24.ru/rest/1536/3pllzln0z4a69b23/crm.company.get"
	excelFileName    = "bitrix_data.xlsx"
)

var accountingMap = map[string]string{
	"768": "Аренда офиса",
	"770": "ФОТ",
	"772": "Возврат покупателям",
	"774": "Дивиденды",
	"776": "Командировочные расходы",
	"778": "Коммунальные услуги",
	"780": "Налоги",
	"782": "Маркетинг",
	"784": "Найм",
	"786": "ПО и лицензии для компании",
	"808": "Настройка CRM-системы",
	"810": "Партнерские вознаграждения",
}

type BitrixResponse struct {
	Result []BitrixItem `json:"result"`
}

type BitrixItem struct {
	ID           string            `json:"ID"`
	NAME         string            `json:"NAME"`
	DATE_CREATE  string            `json:"DATE_CREATE"`
	PROPERTY_634 map[string]string `json:"PROPERTY_634"`
	PROPERTY_638 map[string]string `json:"PROPERTY_638"`
	PROPERTY_732 map[string]string `json:"PROPERTY_732"`
	PROPERTY_628 map[string]string `json:"PROPERTY_628"`
	PROPERTY_632 map[string]string `json:"PROPERTY_632"`
	PROPERTY_630 map[string]string `json:"PROPERTY_630"`
	PROPERTY_636 map[string]string `json:"PROPERTY_636"`
}

func FetchBitrixData() []BitrixItem {
	var allResults []BitrixItem
	next := 0
	counts := 0

	for {
		counts++
		log.Printf("Запрос #%d (next = %d)", counts, next)

		requestBody, _ := json.Marshal(map[string]interface{}{
			"IBLOCK_TYPE_ID": "bitrix_processes",
			"IBLOCK_ID":      "108",
			"start":          next,
		})

		resp, err := http.Post(bitrixAPI, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatalf("Ошибка выполнения запроса: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Ошибка чтения ответа: %v", err)
		}

		log.Printf("Размер ответа: %d байт", len(body))

		var bitrixResponse struct {
			Result []BitrixItem `json:"result"`
			Next   *int         `json:"next"`
		}

		if err := json.Unmarshal(body, &bitrixResponse); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			log.Printf("Тело ответа: %s", string(body))
			continue
		}

		log.Printf("Получено записей: %d", len(bitrixResponse.Result))
		allResults = append(allResults, bitrixResponse.Result...)

		if bitrixResponse.Next == nil {
			log.Println("Данные закончились, выходим из цикла")
			break
		}
		next = *bitrixResponse.Next
	}

	log.Printf("Всего загружено записей: %d", len(allResults))
	return allResults
}

func CreateExcel(data []BitrixItem) {
	f := excelize.NewFile()
	headers := []string{"Приход / Расход", "За что", "Дата операции", "Дата учета", "Сумма",
		"Поступления (от кого)", "Статья учета", "Комментарий", "Дата создания", "ID"}

	sheetName := "Bitrix Data"
	f.SetSheetName("Sheet1", sheetName)

	// Записываем заголовки
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Записываем данные
	for rowIndex, item := range data {
		rowNum := rowIndex + 2
		trans := map[string]string{"790": "Расход", "788": "Приход"}
		incomeExpense := trans[getValue(item.PROPERTY_634)]

		companyOrDealTitle := ""
		for _, value := range item.PROPERTY_632 {
			if strings.HasPrefix(value, "CO_") {
				companyOrDealTitle = fetchTitle(bitrixCompanyAPI, strings.TrimPrefix(value, "CO_"))
			} else if strings.HasPrefix(value, "D_") {
				companyOrDealTitle = fetchTitle(bitrixDealAPI, strings.TrimPrefix(value, "D_"))
			}
		}

		accounting := accountingMap[getValue(item.PROPERTY_630)]
		sumStr := getValue(item.PROPERTY_628)
		sumFloat, _ := strconv.ParseFloat(sumStr, 64)

		values := []interface{}{
			incomeExpense, item.NAME, getValue(item.PROPERTY_638), getValue(item.PROPERTY_732),
			sumFloat, companyOrDealTitle, accounting, getValue(item.PROPERTY_636), item.DATE_CREATE, item.ID,
		}

		for colIndex, value := range values {
			cell := fmt.Sprintf("%c%d", 'A'+colIndex, rowNum)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	if err := f.SaveAs(excelFileName); err != nil {
		log.Fatalf("Ошибка сохранения Excel: %v", err)
	}
	fmt.Println("Excel-файл успешно сохранён:", excelFileName)
}

func getValue(prop map[string]string) string {
	for _, v := range prop {
		return v
	}
	return ""
}

// fetchTitle с увеличенной задержкой между повторными попытками
func fetchTitle(url, id string) string {
	requestBody, _ := json.Marshal(map[string]string{"id": id})

	var response struct {
		Result struct {
			TITLE string `json:"TITLE"`
		} `json:"result"`
		Error string `json:"error"`
	}

	log.Printf("Запрос fetchs Title: URL=%s, ID=%s", url, id)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Ошибка запроса к Bitrix24: %v", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	log.Printf("Ответ fetchTitle для ID=%s: %s", id, string(body))

	time.Sleep(500 * time.Millisecond) // Ограничение: 2 запроса в секунду

	return response.Result.TITLE
}
