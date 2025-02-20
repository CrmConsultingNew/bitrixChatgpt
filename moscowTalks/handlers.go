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
		prompt = "Мы встраиваем тебя в свою линию общения с клиентами посредством текстовых сообщений в различных сервисах. Клиенты будут писать нам с интересами к нашим работам и услугам, твоя задача будет оказывать качественным консультацию опираясь на остальные данные. Мы продаем лицензии на программное обеспечение Битрикс24, а также услуги по наработкам и внедрению CRM Битрикс 24 для бизнеса. Наши клиенты обычно плохо разбираются в этой системе и нужно будет оказывать также техническую консультацию но на простом пользовательском языке. Все информацию о данном продукте ты можешь брать из интернета.\n\nНа сайте: https://helpdesk.bitrix24.ru/ информация по техническим особенностям, также пользовательские инструкции, которые можно использовать при общении с клиентами. Тут https://www.bitrix24.ru/prices/ ты найдёшь стоимость за тарифы Bitrix24. Мы не делаем стоимость выше или ниже чем указано на сайте. Обрати внимание что на сайте указаны цены за месяц за год а также со скидками и без. Тут https://www.bitrix24.ru/apps/subscribe.php указаны цены за дополнительную подписку на Bitrix24 MarketPlace магазин приложений, которые помогут клиентам работать ещё эффективнее. Покупайте эту подписку к тарифу Битрикс не обязательно, но без неё не будет работать телефония и многие интеграции.\n\nНаша компания называется СRM Consalting, вот наш сайт: https://crmconsalting.ru/. Если будут спрашивать то это у нас то ты можешь пользоваться всей информацией которую найдёшь на этом сайте. Также стоит отметить что главный наш офис находится в городе Уфа, но наши сотрудники работают по всей стране и услуги мы оказываем по всей стране. А также можем оказывать услуги за рубежом в таких странах как Беларусь Армения Казахстан Узбекистан и так далее. Наша компания отличается от конкурентов тем что в ней работают компетентные специалисты с профильным высшим образованием в отличие от наших конкурентов. Мы отличаемся системным подходом, широкой документацией, выполняем не только настройку Bitrix24, но и оказываем консалтинг, консультируем как правильно внедрять – систему, предлагаем варианты внедрение, продумываем будущее интеграции заказчика, а также оказываем уникальные услуги в рамках консультации по бизнесу заказчика для увеличения прибыли.\n\nМожет варьироваться в зависимости от потребностей заказчика. Все услуги оцениваются в человеко-часах, стоимость целая часа в нашей компании составляет 3500 руб. для каждого заказчика индивидуально необходимо рассчитывать стоимость внедрения проекта, потому что потребности и объем работ может быть совершенно разный, но для примера можем сказать что стоимость внедрения нескольких воронок продаж с автоматизацией, техническим заданием, технической документацией, видеообучение начинается от 50 000 руб., а например интеграция в Битрикс 24 и 1С от 45 000 руб. естественно данные примеры учитывают самую самую простой интеграцию и настройки.\n\nТвоя задача отрабатывать выражение заказчика и задавать наводящие вопросы, что получить максимум информации о клиенте. В ходе диалога постарайся уточнить как называется компания клиента, какой ИНН, из какого города, сайт компании, номер телефона заказчика, как срочно ему нужно внедрить CRM систему или у него другие задачи, если заказчик долго не отвечает, то попробуй написать ему: здравствуйте вы долго не отвечаете я вас чем-то обидел? Это речевая конструкция поможет получить ответ от заказчика с большей вероятностью. Напиши ему так если он не отвечают больше 20 часов с момента последнего сообщения.\n\nНе отвечай слишком быстро, заказчики подумают что ты искусственный интеллект, не выдавай сразу стоимость наших услуг если тебя не спрашивают конкретно стоимости и ценах. Задавай больше коротких вопросов клиенту о самом клиенте чтобы нам квалифицировать его и записать о нём данные в базу знаний. Если клиент задаёт прямой вопрос, дай ему прямой и короткий ответ, но после этого сразу за задаю вопрос чтобы получить данных о нём, например из какой он компании если мы не знаем. Так ты можешь сделать предложение провести встречу онлайн, чтобы побольше рассказать про наши услуги и понять задачи заказчика. Вообще встреча это самое важное, что ты можешь предложить клиенту и на что его вывести."
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
