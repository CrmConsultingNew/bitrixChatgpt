package torgi_new

/*func CheckIfValidOrder(url string) bool {
	// Загружаем HTML-страницу закупки
	htmlContent, err := GetInnPage(url)
	if err != nil {
		log.Printf("Ошибка загрузки закупки %s: %s\n", url, err)
		return false
	}

	// Парсим страницу с помощью goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("Ошибка парсинга HTML закупки %s: %s\n", url, err)
		return false
	}

	// Ищем текст "Электронная Торговая Площадка Торги-Онлайн"
	doc.Find(".common-text__value").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "Электронная Торговая Площадка Торги-Онлайн") {
			log.Printf("⚠️ Закупка %s отклонена (найдено 'Электронная Торговая Площадка Торги-Онлайн')\n", url)
		}
	})

	// Возвращаем `false`, если закупка невалидна, иначе `true`
	return !doc.Find(".common-text__value").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(strings.TrimSpace(s.Text()), "Электронная Торговая Площадка Торги-Онлайн")
	}).NodesAreEmpty()
}

func GetInnPage(url string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Ошибка создания запроса: %w", err)
	}

	// Устанавливаем User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	// Читаем содержимое страницы
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Ошибка чтения ответа: %w", err)
	}

	return string(body), nil
}
*/
