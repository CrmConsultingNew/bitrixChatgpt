package evroangar

import (
	"io"
	"log"
	"net/http"
)

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Webhook received, processing...")

	// Считываем тело запроса
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Логируем сырые данные запроса
	log.Println("Raw request body:", string(bs))
}
