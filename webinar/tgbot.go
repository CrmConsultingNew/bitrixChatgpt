package webinar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	botToken = "780504069:AAH7Ld_hobbvEkCZi8fpdKUIEXirpG4raCQ"
	dataFile = "telegramContactsCrmConsulting.json"
)

// Структура клиента
type Client struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Глобальная мапа с уникальными пользователями
var uniqueClients map[int64]Client

// Запуск бота
func StartTgBot() {
	log.Println("Starting Telegram Bot")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(Handler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %v", err)
	}

	// Загружаем клиентов
	loadUniqueClients()

	// Запускаем планировщик сообщений
	go ScheduleMessages(b)

	b.Start(ctx)
	log.Println("b.ID():", b.ID())
}

// Загрузка уникальных клиентов в память
func loadUniqueClients() {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		log.Printf("Файл не найден или ошибка чтения: %v", err)
		return
	}

	var clients []Client
	if err := json.Unmarshal(file, &clients); err != nil {
		log.Fatalf("Ошибка при разборе JSON: %v", err)
	}

	uniqueClients = make(map[int64]Client)
	for _, client := range clients {
		uniqueClients[client.UserID] = client
	}

	log.Printf("Загружено %d уникальных пользователей", len(uniqueClients))
}

// Сохранение клиента в JSON, если его еще нет
func saveClientData(client Client) {
	clients := loadClients()

	// Добавляем нового клиента
	clients = append(clients, client)

	// Перезаписываем JSON-файл с массивом (валидный JSON)
	file, err := os.OpenFile(dataFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Красивый JSON с отступами
	err = encoder.Encode(clients)
	if err != nil {
		log.Printf("Ошибка при записи в файл: %v", err)
	}
}

func loadClients() []Client {
	var clients []Client

	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Client{} // Файл еще не создан
		}
		log.Printf("Ошибка при чтении файла: %v", err)
		return []Client{}
	}

	err = json.Unmarshal(file, &clients)
	if err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		return []Client{}
	}

	return clients
}

// Универсальная функция отправки сообщений
func sendMessageToClients(b *bot.Bot, messageTemplate string) {
	ctx := context.Background()
	for userID, client := range uniqueClients {
		var message string
		if strings.Contains(messageTemplate, "%s") {
			message = fmt.Sprintf(messageTemplate, client.FirstName)
		} else {
			message = messageTemplate
		}

		log.Printf("Отправка сообщения пользователю %d: %s", userID, message)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    userID,
			Text:      message,
			ParseMode: "MarkdownV2",
		})
		if err != nil {
			log.Printf("Ошибка при отправке сообщения пользователю %d: %v", userID, err)
		} else {
			log.Printf("Сообщение успешно отправлено пользователю %d", userID)
		}
	}
}

// Планировщик сообщений
func ScheduleMessages(b *bot.Bot) {
	for {
		now := time.Now().Format("15:04") // Время в формате ЧЧ:ММ

		switch now {
		case "09:00":
			sendMessageToClients(b, `*Вебинар через 2 часа\! Как CRM \+ AI увеличивают прибыль\?*

%s, доброе утро\!

На связи команда CRM Consulting, через 2 часа мы разберём, как автоматизировать продажи, перестать терять клиентов и увеличить прибыль с помощью CRM и искусственного интеллекта\.

Вот что вы узнаете на вебинаре:
• Как CRM \+ AI помогают не терять до 40%% клиентов – автоматизация заявок, напоминания, персонализированные follow\-up\'ы\.
• Почему CRM – это не просто база данных, а инструмент роста – реальные кейсы, как бизнесы увеличили выручку на 30\-50%%\.
• Как построить отдел продаж, который работает без хаоса – CRM поможет внедрить чёткие процессы\.

Подключайтесь вовремя, чтобы не пропустить самое важное\!
👉 [Присоединяйтесь к вебинару здесь](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)

Скоро увидимся в эфире\!
Команда CRM Consulting`)
		case "10:50":
			sendMessageToClients(b, `*Осталось 10 минут и начинаем\. Подключайтесь*

%s, уже через 10 минут разберем, как CRM \+ AI могут закрыть ключевые проблемы бизнеса и помочь вам зарабатывать больше\!

Вам будет полезно, если\:

• *Менеджеры теряют клиентов* – заявки обрабатываются с опозданием, клиенты уходят к конкурентам\.
• *Нет прозрачности в продажах* – сложно понять, на каком этапе сделки, кто выполняет задачи и где узкие места\.
• *Слабый контроль за командой* – не хватает данных, чтобы объективно оценить работу менеджеров\.
• *Ручная работа замедляет бизнес* – отчёты, напоминания, письма менеджеры делают вручную вместо продаж\.
• CRM есть, но работает вполсилы или вы задумываетесь о внедрении\.

🔥 Сегодня разберём, как CRM \+ AI решают эти проблемы и превращают хаос в систему, которая реально продает\!

[Ваша ссылка на подключение](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)`)
		case "11:00":
			sendMessageToClients(b, `*3\.\.\.2\.\.\.1\.\.\. Мы в эфире\! Подключайтесь прямо сейчас\!*

Мы уже начали вебинар\! *Как CRM \+ AI увеличивают продажи, сокращают потери клиентов и систематизируют бизнес\?* Разберем прямо сейчас\!

Обсудим без воды: как перестать терять клиентов и увеличить прибыль с CRM, какие процессы можно автоматизировать уже сейчас, реальные кейсы компаний,
которые уже используют наши клиенты\.

Присоединяйтесь 👉 [Ссылка на вебинар](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)

Не упустите, то что для вас важно \- присоединяйтесь\!`)
		case "11:15":
			sendMessageToClients(b, `*Как подружить CRM и искусственный интеллект\?*

Михаил, бизнес\-аналитик CRM Consulting рассказывает в эфире прямо сейчас, сэкономить и увеличить продажи в бизнесе с помощью искусственного интеллекта\.

Заходите 👉 [Ссылка](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)
`)
		case "11:30":
			sendMessageToClients(b, `*🔥 Снизили затраты на Call\-центр в 10 раз\! С помощью ИИ*

Сегодня на вебинаре мы рассматриваем реальные истории внедрения ИИ \+ CRM\.

Подключайтесь, чтобы узнать, как компания отказалась от Call\-центра на аутсорсе благодаря искусственному интеллекту\.

Ждем вас: 👉 [Ссылка на вебинар](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)
`)
		}

		time.Sleep(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	}
}
