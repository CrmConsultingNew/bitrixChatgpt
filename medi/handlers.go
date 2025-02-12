package medi

import (
	"fmt"
	"net/http"
	"time"
)

// StartMedi - основной обработчик
func StartMedi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	GetContactsListWithCustomFieldsBirthdate()

	contacts := ReadContactsJsonAndGetClientContactPhone()
	if len(contacts) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Нет контактов с ДР на сегодня."))
		return
	}

	// Ждем до 08:00 по Москве
	//waitUntilEightAM()

	// Отправляем поздравления
	for _, contact := range contacts {
		GlobalTextMessageToClient = fmt.Sprintf(
			"Уважаемый %s %s %s. Компания Меди поздравляет с днем рождения!",
			contact["last_name"], contact["name"], contact["second_name"],
		)
		SendMessageToClient(contact["phone"])
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Контакты обработаны... StartMedi func"))
}

// Ожидание 08:00 по московскому времени
func waitUntilEightAM() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Ошибка загрузки часового пояса:", err)
		return
	}

	for {
		currentTime := time.Now().In(loc)
		if currentTime.Hour() == 8 && currentTime.Minute() == 0 {
			break // Время 08:00, выходим из цикла
		}
		time.Sleep(1 * time.Minute) // Проверяем каждую минуту
	}
}
