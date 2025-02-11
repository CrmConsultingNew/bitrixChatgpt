package metrologiya

import (
	"github.com/go-co-op/gocron"
	"log"
	"net/http"
	"time"
)

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

// StartScheduler запускает планировщик задач
func StartScheduler() {
	scheduler := gocron.NewScheduler(time.Local)

	// torgi_online_cron
	scheduler.Every(1).Day().At("17:00").Do(GenerateDailyReportForTorgiOnline)

	// Ежедневный отчет (каждый день в 09:00)
	scheduler.Every(1).Day().At("09:00").Do(GenerateDailyReport)

	// Еженедельный отчет (каждую среду в 09:00)
	scheduler.Every(1).Wednesday().At("09:00").Do(GenerateWeeklyReport)

	// Ежемесячный отчет (каждое 15 число каждого месяца в 09:00)
	scheduler.Every(1).Month(15).At("09:00").Do(GenerateMonthlyReport)

	log.Println("Планировщик задач запущен...")
	scheduler.StartAsync()
}

// GenerateDailyReport отправляет вебхук для ежедневного отчета
func GenerateDailyReportForTorgiOnline() {
	log.Println("Отправка вебхука для ежедневного отчета torgi-online...")
	sendWebhook("https://crmconsulting-api.ru/api/start_torgi")
}

// GenerateDailyReport отправляет вебхук для ежедневного отчета
func GenerateDailyReport() {
	log.Println("Отправка вебхука для ежедневного отчета...")
	sendWebhook("https://crmconsulting-api.ru/api/request")
}

// GenerateWeeklyReport отправляет вебхук для еженедельного отчета
func GenerateWeeklyReport() {
	log.Println("Отправка вебхука для еженедельного отчета...")
	sendWebhook("https://crmconsulting-api.ru/api/weekly")
}

// GenerateMonthlyReport отправляет вебхук для ежемесячного отчета
func GenerateMonthlyReport() {
	log.Println("Отправка вебхука для ежемесячного отчета...")
	sendWebhook("https://crmconsulting-api.ru/api/monthly")
}

func CloseAllBizProcesses() {
	log.Println("Отправка вебхука для ежемесячного отчета...")
	sendWebhook("https://crmconsulting-api.ru/api/close_bizproc")
}
