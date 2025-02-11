package reports

import "time"

const (
	DEBUG                  = true
	LOG_DIR                = "" // Пустая строка для директории логов
	LOG_FILE               = "log.txt"
	B24_HOST               = "torgi-online.bitrix24.ru"
	B24_USER_ID            = 88
	B24_WEBHOOK            = "672m3xnegwu7bosa"
	B24SLEEP               = 500 * time.Millisecond // В Go вместо микросекунд используется time.Duration
	B24_LEAD_REGDATE_FIELD = "UF_CRM_1614057621"
	B24_LEAD_KPDATE_FIELD  = "UF_CRM_1614648126"
	B24_LEAD_KP_STATUS_ID  = "11"
	B24_TASK_SUCCTAG       = "ок"
	B24_MANAGER_DEPARTMENT = 5
	B24_IT_DEPARTMENT      = 10
)

// Функция для проверки отладки
func IsDebug() bool {
	return DEBUG
}
