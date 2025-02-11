package torgi_online

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func TestTorgi() {

	noticeInfoID := "17834693" // Пример ID извещения
	processOrderInfo(noticeInfoID)
}

func FetchHtmlOld(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 HTTP status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

func processOrderInfo(noticeInfoID string) {
	// Формируем URL для запроса
	url := fmt.Sprintf("https://zakupki.gov.ru/epz/order/notice/notice223/common-info.html?noticeInfoId=%s", noticeInfoID)

	log.Printf("Загружаемый URL: %s", url)

	// Загружаем HTML-страницу
	doc, err := FetchHtmlOld(url)
	if err != nil {
		log.Printf("Ошибка загрузки HTML для noticeInfoID %s: %v", noticeInfoID, err)
		return
	}

	// Извлекаем значение элемента #7
	element := doc.Find(".common-text__value").Eq(6) // Индекс 6 для седьмого элемента (нумерация начинается с 0)
	text := strings.TrimSpace(element.Text())

	log.Printf("Значение элемента #7: %s", text)

	// Дополнительная логика, если требуется
	if strings.Contains(text, "Электронная Торговая Площадка Торги-Онлайн") {
		log.Printf("Процесс прерван для noticeInfoID %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", noticeInfoID)
		return
	}

	log.Printf("Обработка завершена для noticeInfoID %s.", noticeInfoID)
}

func cleanNumberString(input string) string {
	// Убираем всё, кроме цифр
	return strings.Join(strings.FieldsFunc(input, func(r rune) bool {
		return r < '0' || r > '9'
	}), "")
}

func processCompanies(compInn string, todayNew string, allOrdersHrefs map[string]struct{}) {
	// Формируем URL для запроса
	url := fmt.Sprintf("https://zakupki.gov.ru/epz/order/extendedsearch/results.html?searchString=%s&morphology=on&search-filter=\u0414\u0410\u0422\u0415+\u0420\u0410\u0417\u041c\u0415\u0429\u0415\u041d\u0418\u042f&pageNumber=1&sortDirection=false&recordsPerPage=_10&showLotsInfoHidden=false&sortBy=UPDATE_DATE&fz223=on&af=on&currencyIdGeneral=-1&publishDateFrom=%s", compInn, todayNew)

	log.Printf("Загружаемый URL: %s", url)

	// Загружаем HTML-страницу
	doc, err := FetchHtmlOld(url)
	if err != nil {
		log.Printf("Ошибка загрузки HTML для ИНН %s: %v", compInn, err)
		return
	}

	// Извлекаем все значения div с классом common-text__value
	log.Println("Извлекаем все значения элементов <div class='common-text__value'>:")
	doc.Find(".r").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		log.Printf("Элемент #%d: %s", i+1, text)
	})

	// Проверяем все элементы .common-text__value на наличие искомого текста
	log.Println("Ищем элемент с текстом 'Электронная Торговая Площадка Торги-Онлайн'")
	found := false
	doc.Find(".common-text__value").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "Электронная Торговая Площадка Торги-Онлайн") {
			log.Printf("Процесс прерван для ИНН %s: найден текст 'Электронная Торговая Площадка Торги-Онлайн'.", compInn)
			found = true
			return
		}
	})

	if found {
		return
	}

	// Обработка новых заказов
	newOrdersText := doc.Find(".search-results__total").First().Text()
	cleanedOrdersText := cleanNumberString(newOrdersText)

	newOrders, err := strconv.Atoi(cleanedOrdersText)
	if err != nil {
		log.Printf("Ошибка преобразования числа новых заказов для ИНН %s: %v (original text: %q)", compInn, err, newOrdersText)
		newOrders = 0
	}

	log.Printf("Найдено %d новых заказов для ИНН %s", newOrders, compInn)

	if newOrders > 0 {
		// Извлечение дополнительной информации
		compHref, exists := doc.Find(".registry-entry__header-mid__number a").Attr("href")
		if !exists || compHref == "" {
			log.Printf("Ссылка на заказ не найдена для ИНН %s", compInn)
			return
		}
		if !strings.Contains(compHref, "torgi-online.gov.ru") {
			compHref = "https://zakupki.gov.ru/" + compHref
		}

		dateStart := doc.Find(".data-block__value").Eq(0).Text()
		dateUpdate := doc.Find(".data-block__value").Eq(1).Text()
		dateFinish := doc.Find(".data-block__value").Eq(2).Text()

		if _, exists := allOrdersHrefs[compHref]; !exists {
			log.Printf("Новый заказ: ИНН %s, Размещено/Обновлено/Окончание - %s/%s/%s, Госзакупка - %s", compInn, dateStart, dateUpdate, dateFinish, compHref)
			allOrdersHrefs[compHref] = struct{}{}
		}
	}
}
