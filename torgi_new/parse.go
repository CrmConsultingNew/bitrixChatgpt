package torgi_new

import (
	torgi_online "bitrix_app/torgi-online"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var CompCount int

// extractNumberFromText извлекает число из строки с текстом
func extractNumberFromText(text string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(text)
	num, err := strconv.Atoi(match)
	if err != nil {
		return 0
	}
	return num
}

// IsValidOrder Проверяет, не относится ли закупка к "Электронная Торговая Площадка Торги-Онлайн"
func IsValidOrder(url string) bool {
	log.Printf("🔍 Проверяем закупку: %s", url)

	// Загружаем HTML страницы закупки
	htmlContent, err := GetInn(url, "")
	if err != nil {
		log.Printf("❌ Ошибка загрузки закупки %s: %s\n", url, err)
		return false
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("❌ Ошибка парсинга HTML закупки %s: %s\n", url, err)
		return false
	}

	found := false

	// Загружаем HTML-страницу
	oldDoc, err := torgi_online.FetchHtmlOld(url)
	if err != nil {
		log.Printf("Ошибка загрузки HTML для url %s: %v", url, err)
		return false
	}

	// Извлекаем значение элемента #7
	elementOne := oldDoc.Find(".common-text__value").Eq(1) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textOne := strings.TrimSpace(elementOne.Text())
	elementTwo := oldDoc.Find(".common-text__value").Eq(2) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textTwo := strings.TrimSpace(elementTwo.Text())
	elementThree := oldDoc.Find(".common-text__value").Eq(3) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textThree := strings.TrimSpace(elementThree.Text())
	elementFour := oldDoc.Find(".common-text__value").Eq(4) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textFour := strings.TrimSpace(elementFour.Text())
	elementFive := oldDoc.Find(".common-text__value").Eq(5) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textFive := strings.TrimSpace(elementFive.Text())
	elementSix := oldDoc.Find(".common-text__value").Eq(6) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textSix := strings.TrimSpace(elementSix.Text())
	elementSeven := oldDoc.Find(".common-text__value").Eq(7) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textSeven := strings.TrimSpace(elementSeven.Text())
	elementEight := oldDoc.Find(".common-text__value").Eq(8) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	textEight := strings.TrimSpace(elementEight.Text())

	//log.Printf("Значение элемента #7: %s", text)

	// Дополнительная логика, если требуется
	if strings.Contains(textOne, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	if strings.Contains(textTwo, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	if strings.Contains(textThree, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	if strings.Contains(textFour, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	if strings.Contains(textFive, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	if strings.Contains(textSix, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	if strings.Contains(textSeven, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	if strings.Contains(textEight, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для url %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", url)
		return false
	}
	// Находим заголовки `.common-text__title` и проверяем их содержимое
	doc.Find(".common-text__title").Each(func(i int, s *goquery.Selection) {
		log.Printf("🔎 Заголовок #%d: %s", i+1, s.Text())
		text := strings.TrimSpace(s.Text())

		// Если нашли "Наименование электронной площадки"
		if strings.Contains(text, "Наименование электронной площадки") {
			// Следующий `.common-text__value` содержит название площадки
			value := s.Parent().Find(".common-text__value").First().Text()
			value = strings.TrimSpace(value)

			log.Printf("📝 Найдено название площадки: '%s'", value) // Отладка

			if strings.Contains(value, "Электронная Торговая Площадка Торги-Онлайн") {
				log.Printf("⚠️ Закупка %s исключена (ЭТП Торги-Онлайн)\n", url)
				found = true
			}
		}
	})

	return !found
}

// ProcessCompanies - обрабатывает все компании из списка
func ProcessCompanies(arrayAllComp []map[string]interface{}, allOrdersHrefs []string, arUsers map[int]string, logTraceID string) ([]string, map[int]string) {
	todayNew := time.Now().AddDate(0, 0, -3).Format("02.01.2006")

	newCompaniesListArr := make(map[int]string)
	arrKey := 0

	processedCount := 0

	for _, oneComp := range arrayAllComp {

		COMPInn, ok := oneComp["UF_INN"].(string)
		if !ok || COMPInn == "" {
			continue
		}

		COMPIDStr, ok := oneComp["ID"].(string)
		if !ok {
			COMPIDStr = fmt.Sprintf("%.0f", oneComp["ID"].(float64))
		}
		COMPID, _ := strconv.Atoi(COMPIDStr)

		COMPUserIDStr, ok := oneComp["ASSIGNED_BY_ID"].(string)
		if !ok {
			COMPUserIDStr = fmt.Sprintf("%.0f", oneComp["ASSIGNED_BY_ID"].(float64))
		}
		COMPUserID, _ := strconv.Atoi(COMPUserIDStr)

		COMPName, _ := oneComp["TITLE"].(string)

		log.Printf("[%s] Обработка ИНН: %s (№%d)\n", logTraceID, COMPInn, processedCount+1)

		htmlContent, err := GetInn(COMPInn, todayNew)
		if err != nil {
			log.Printf("[%s] Ошибка загрузки HTML для ИНН %s: %s\n", logTraceID, COMPInn, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			log.Printf("[%s] Ошибка парсинга HTML для ИНН %s: %s\n", logTraceID, COMPInn, err)
			continue
		}

		totalText := doc.Find(".search-results__total").First().Text()
		totalText = strings.TrimSpace(totalText)

		log.Printf("[%s] Найден текст в .search-results__total для ИНН %s: '%s'\n", logTraceID, COMPInn, totalText)

		newOrders := extractNumberFromText(totalText)
		newOrdersParsed, compHref, dateStart, dateUpdate, dateFinish := ParseHTML(doc)

		if newOrdersParsed > newOrders {
			newOrders = newOrdersParsed
		}

		if newOrders == 0 || compHref == "" {
			log.Printf("[%s] Нет новых закупок для ИНН %s\n", logTraceID, COMPInn)
			continue
		}

		if Contains(allOrdersHrefs, compHref) {
			log.Printf("[%s] Закупка уже обработана: %s\n", logTraceID, compHref)
			continue
		}

		// Проверяем, относится ли закупка к "ЭТП Торги-Онлайн"
		if !IsValidOrder(compHref) {
			log.Printf("[%s] ❌ Закупка %s исключена (ЭТП Торги-Онлайн)\n", logTraceID, compHref)
			continue
		}

		manager, exists := arUsers[COMPUserID]
		if !exists {
			manager = "Неизвестный менеджер"
		}

		entry := fmt.Sprintf("%d) %s; ИНН - %s; Размещено/Обновлено/Окончание - %s / %s / %s; Компания - https://torgi-crm.online/crm/company/details/%d/; Госзакупка - %s - %s\n",
			CompCount+1, COMPName, COMPInn, dateStart, dateUpdate, dateFinish, int(COMPID), compHref, manager)

		if CompCount%40 == 0 && CompCount != 0 {
			arrKey++
		}

		log.Printf("[%s] ✅ Найдена новая закупка для ИНН %s: %s\n", logTraceID, COMPInn, compHref)

		newCompaniesListArr[arrKey] += entry
		allOrdersHrefs = append(allOrdersHrefs, compHref)
		CompCount++
		processedCount++
	}

	log.Printf("[%s] Всего новых закупок найдено: %d\n", logTraceID, CompCount)
	return allOrdersHrefs, newCompaniesListArr
}

func ParseHTML(doc *goquery.Document) (int, string, string, string, string) {
	var (
		newOrders  int
		compHref   string
		dateStart  string
		dateUpdate string
		dateFinish string
	)

	// Извлекаем значение элемента #7
	element := doc.Find(".common-text__value").Eq(6) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	text := strings.TrimSpace(element.Text())

	log.Printf("Значение элемента #7: %s", text)

	// Дополнительная логика, если требуется
	if strings.Contains(text, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван, найден текст 'Электронная Торговая Площадка Торги-Онлайн'")
		return 0, "", "", "", ""
	}

	// Extract total number of orders
	totalText := doc.Find(".search-results__total").First().Text()
	if num, err := strconv.Atoi(strings.TrimSpace(totalText)); err == nil {
		newOrders = num
	}

	// Extract procurement link
	doc.Find(".registry-entry__header-mid__number a").First().Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			compHref = href
			if !strings.HasPrefix(compHref, "https://") {
				compHref = "https://zakupki.gov.ru" + compHref
			}
		}
	})

	// Extract date values
	doc.Find(".data-block__value").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			dateStart = s.Text()
		case 1:
			dateUpdate = s.Text()
		case 2:
			dateFinish = s.Text()
		}
	})

	return newOrders, compHref, dateStart, dateUpdate, dateFinish
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// fetchHTML Deprecated?
func fetchHTML(inn, date string) *goquery.Document {
	time.Sleep(500 * time.Millisecond) // Delay to prevent server blocking

	url := fmt.Sprintf("https://zakupki.gov.ru/epz/order/extendedsearch/results.html?searchString=%s&morphology=on&search-filter=%D0%94%D0%B0%D1%82%D0%B5+%D1%80%D0%B0%D0%B7%D0%BC%D0%B5%D1%89%D0%B5%D0%BD%D0%B8%D1%8F&pageNumber=1&sortDirection=false&recordsPerPage=_10&showLotsInfoHidden=false&sortBy=UPDATE_DATE&fz223=on&af=on&currencyIdGeneral=-1&publishDateFrom=%s", inn, date)

	// Выполняем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch URL: %s, error: %s\n", url, err)
		return nil
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected HTTP status: %d for URL: %s\n", resp.StatusCode, url)
		return nil
	}

	// Создаем объект Document из тела ответа
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse HTML for INN %s: %s\n", inn, err)
		return nil
	}

	return doc
}
