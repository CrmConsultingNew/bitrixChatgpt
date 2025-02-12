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

	// Ждем до 08:00 по Москве (если нужно)
	// waitUntilEightAM()

	// Сообщение для отправки
	message := "Скоро Ваш день рождения!\n" +
		"Для каждого из нас - это особая дата!\n" +
		"И, конечно же, никакой День рождения не может обойтись без приятных сюрпризов.\n" +
		"А вот и они! Примите от нас искренние поздравления и подарок, " +
		"дополнительная 10 % скидка на наши услуги.\n" +
		"*Скидка действует в течение 5 (пяти) дней до и после даты рождения."

	// Отправляем сообщение каждому контакту
	for _, contact := range contacts {
		SendMessageToClient(contact["phone"], message)
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
