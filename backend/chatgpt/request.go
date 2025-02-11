package chatgpt

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func RequestFromVue(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Message string `json:"message"`
		DealId  string `json:"dealId"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		log.Println("Error unmarshalling request body:", err)
		http.Error(w, "Error in request format", http.StatusBadRequest)
		return
	}

	log.Println("req_at_Vue_for_chatgpt:", requestData.Message)

	// Use the extracted message in your SendRequest function
	responseText := SendRequest(requestData.Message)

	// Write the SendRequest response back to the frontend
	w.Write([]byte(responseText))
}
