package medi

import (
	"github.com/go-co-op/gocron"
	"log"
	"net/http"
	"time"
)

// StartMediScheduler запускает второй планировщик задач
func StartMediScheduler() {
	scheduler := gocron.NewScheduler(time.Local)

	// Новый CRON — проверка дат контактов каждый день в 01:00
	scheduler.Every(1).Day().At("22:29").Do(CheckEveryDayContactsDate)

	log.Println("Планировщик задач StartMediScheduler запущен...")
	scheduler.StartAsync()
}

// CheckEveryDayContactsDate проверяет даты контактов каждую ночь в 01:00
func CheckEveryDayContactsDate() {
	log.Println("Отправка вебхука для проверки дат контактов...")
	sendWebhook("https://crmconsulting-api.ru/api/medi_birthdate")
}

// sendWebhook отправляет GET-запрос на указанный endpoint
func sendWebhook(endpoint string) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("Ошибка создания запроса для %s: %v", endpoint, err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка выполнения запроса для %s: %v", endpoint, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Успешно отправлен вебхук на %s, статус: %s", endpoint, resp.Status)
}
