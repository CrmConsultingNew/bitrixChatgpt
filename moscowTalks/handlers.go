package moscowTalks

import (
	"bitrix_app/OpenAI"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

// Хранилище состояния приветствий
var greetingCache = struct {
	mu             sync.Mutex
	greetedDialogs map[string]bool
}{
	greetedDialogs: make(map[string]bool),
}

type MessageData struct {
	Event                string
	EventHandlerID       string
	DialogID             string
	ChatEntityID         string
	Domain               string
	ChatEntityData1      string
	ChatEntityData2      string
	Message              string
	MessageID            string
	UserName             string
	UserFirstName        string
	UserLastName         string
	UserGender           string
	UserIsBot            string
	UserIsConnector      string
	UserIsNetwork        string
	UserIsExtranet       string
	AuthClientEndpoint   string
	AuthServerEndpoint   string
	AuthMemberId         string
	AuthApplicationToken string
	BotCode              string
	BotID                string
}

// Функция для динамического извлечения BOT_ID и BOT_CODE
func extractBotData(values url.Values) (string, string) {
	botID := ""
	botCode := ""

	for key, val := range values {
		if strings.HasPrefix(key, "data[BOT][") {
			if strings.HasSuffix(key, "][BOT_ID]") {
				botID = val[0]
			}
			if strings.HasSuffix(key, "][BOT_CODE]") {
				botCode = val[0]
			}
		}
	}

	return botID, botCode
}

func HandleMessageFromOpenline(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling message was called")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}

	log.Println("string(body):::", string(body))

	decoded, err := url.QueryUnescape(string(body))
	if err != nil {
		log.Println("Error decoding URL:", err)
		http.Error(w, "Error decoding URL", http.StatusInternalServerError)
		return
	}

	values, err := url.ParseQuery(decoded)
	if err != nil {
		log.Println("Error parsing query:", err)
		http.Error(w, "Error parsing query", http.StatusInternalServerError)
		return
	}

	// Динамическое извлечение BOT_ID и BOT_CODE
	botID, botCode := extractBotData(values)

	// Извлекаем значение поля MESSAGE
	rawMessage := values.Get("data[PARAMS][MESSAGE]")
	const prefix = "=== Исходящее сообщение, автор: Телефон ==="
	message := rawMessage
	if strings.Contains(rawMessage, prefix) {
		message = strings.TrimSpace(strings.Replace(rawMessage, prefix, "", 1))
	}

	// Извлекаем интересующие параметры
	data := MessageData{
		Event:                values.Get("event"),
		EventHandlerID:       values.Get("event_handler_id"),
		DialogID:             values.Get("data[PARAMS][CHAT_ID]"),
		ChatEntityID:         values.Get("data[PARAMS][CHAT_ENTITY_ID]"),
		Domain:               values.Get("auth[domain]"),
		ChatEntityData1:      values.Get("data[PARAMS][CHAT_ENTITY_DATA_1]"),
		ChatEntityData2:      values.Get("data[PARAMS][CHAT_ENTITY_DATA_2]"),
		Message:              message,
		MessageID:            values.Get("data[PARAMS][MESSAGE_ID]"),
		UserName:             values.Get("data[USER][NAME]"),
		UserFirstName:        values.Get("data[USER][FIRST_NAME]"),
		UserLastName:         values.Get("data[USER][LAST_NAME]"),
		UserGender:           values.Get("data[USER][GENDER]"),
		UserIsBot:            values.Get("data[USER][IS_BOT]"),
		UserIsConnector:      values.Get("data[USER][IS_CONNECTOR]"),
		UserIsNetwork:        values.Get("data[USER][IS_NETWORK]"),
		UserIsExtranet:       values.Get("data[USER][IS_EXTRANET]"),
		AuthClientEndpoint:   values.Get("auth[client_endpoint]"),
		AuthServerEndpoint:   values.Get("auth[server_endpoint]"),
		AuthMemberId:         values.Get("auth[member_id]"),
		AuthApplicationToken: values.Get("auth[application_token]"),
		BotID:                botID,
		BotCode:              botCode,
	}

	// Логируем данные
	log.Println("DIALOG_ID :", data.DialogID)
	log.Println("BOT_ID :", data.BotID)
	log.Println("BOT_CODE :", data.BotCode)
	log.Println("MESSAGE :", data.Message)

	// Проверка, здоровался ли бот
	greetingCache.mu.Lock()
	alreadyGreeted := greetingCache.greetedDialogs[data.DialogID]
	if !alreadyGreeted {
		// Устанавливаем, что бот уже здоровался в этом диалоге
		greetingCache.greetedDialogs[data.DialogID] = true
	}
	greetingCache.mu.Unlock()

	// Выбор prompt на основе состояния
	prompt := ""
	if alreadyGreeted {
		prompt = "Выяви потребности клиента и задай уточняющие вопросы. Не нужно повторять приветствие."
	} else {
		prompt = "Представь, что ты менеджер по продажам. Тебе нужно узнать имя клиента, когда будет удобно с ним связаться, и выявить потребность клиента. Нельзя говорить что ты ИИ. Тебя зовут Василиса."
	}

	response, err := OpenAI.ChatGptTest(prompt, data.Message)
	if err != nil {
		log.Println("Ошибка при вызове ChatGptTest:", err)
	} else {
		log.Println("Ответ от ChatGPT:", response)
	}

	// Обработка событий
	if data.Event == "ONIMBOTJOINCHAT" {
		log.Println("Ignoring ONIMBOTJOINCHAT event")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Event ignored"))
		return
	}

	if data.UserIsConnector == "N" && data.UserIsExtranet == "N" {
		err := leaveBotFromChat(data.BotID, data.DialogID)
		if err != nil {
			log.Println("Error leaving bot:", err)
		}
		return
	}

	var botClientId string
	switch data.AuthClientEndpoint {
	case "https://crmconsalting.bitrix24.ru/rest/":
		botClientId = os.Getenv("CONSULTING_BOT_CLIENT_ID")
	case "https://b24-iesb30.bitrix24.ru/rest/":
		botClientId = os.Getenv("MY_BOT_CLIENT_ID")
	}
	err = sendMessageToBitrixChatFromChatgpt(data.BotID, data.DialogID, response, botClientId, data.AuthClientEndpoint)
	if err != nil {
		log.Println("Error sending message to chat:", err)
	}
}

func sendMessageToBitrixChatFromChatgpt(BotID, DialogID, message, botClientId, clientEndpoint string) error {
	var webhookData string
	switch clientEndpoint {
	case "https://crmconsalting.bitrix24.ru/rest/":
		webhookData = os.Getenv("CRM_CONSULTING_WEBHOOK")
	case "https://b24-iesb30.bitrix24.ru/rest/":
		webhookData = os.Getenv("MY_WEBHOOK")
	}
	log.Println("webHookData result:", webhookData)
	command := "imbot.message.add"
	webHookUrl := fmt.Sprintf("%s%s", webhookData, command)
	log.Println("webHookUrl result:", webHookUrl)
	requestBody := map[string]interface{}{"BOT_ID": BotID, "DIALOG_ID": "chat" + DialogID, "MESSAGE": message, "CLIENT_ID": botClientId}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("ошибка сериализации тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания HTTP-запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения HTTP-запроса: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

func leaveBotFromChat(BotID, ChatID string) error {
	webHookUrl := "https://b24-iesb30.bitrix24.ru/rest/1/ytdf4fz89gmf7wsp/imbot.chat.leave"
	requestBody := map[string]interface{}{"BOT_ID": BotID, "CHAT_ID": ChatID}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("ошибка сериализации тела запроса: %v", err)
	}

	req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания HTTP-запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения HTTP-запроса: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
