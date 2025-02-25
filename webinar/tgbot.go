package webinar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"os"
	"os/signal"
	"time"
)

//780504069:AAH7Ld_hobbvEkCZi8fpdKUIEXirpG4raCQ

const webinarDate = "2025-02-27"   // Дата вебинара
const testUserID int64 = 947654127 // Тестовый user_id
const botToken = "780504069:AAH7Ld_hobbvEkCZi8fpdKUIEXirpG4raCQ"

// Функция проверки даты и отправки напоминания
func SendReminder(b *bot.Bot) {
	today := time.Now().Format("2006-01-02")
	reminderDate := "2025-02-25" // За два дня до вебинара

	if today == reminderDate {
		message := fmt.Sprintf(`*Напоминание о вебинаре\!*  

Здравствуйте\! Мы напоминаем вам, что уже через 2 дня состоится наш вебинар по внедрению AI в CRM\-системы\.  

📅 *Дата и время:* *%s, 11:00 МСК*  
📌 Ссылка на подключение будет отправлена вам перед мероприятием\.  

Ждем вас на вебинаре\!`, webinarDate)

		ctx := context.Background()
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    testUserID,
			Text:      message,
			ParseMode: "MarkdownV2",
		})

		if err != nil {
			log.Printf("Ошибка при отправке напоминания: %v", err)
		} else {
			log.Println("Напоминание успешно отправлено тестовому пользователю")
		}
	}
}

// Структура для хранения данных клиента
type Client struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Файл для хранения данных клиентов
const dataFile = "telegramContactsCrmConsulting.json"

func StartTgBot() {
	log.Println("Starting Telegram Bot")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	// Запускаем отправку напоминания
	go SendReminder(b)

	b.Start(ctx)
	log.Println("b.ID():", b.ID())
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	// Получаем данные пользователя
	user := update.Message.From
	client := Client{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
	}

	// Логируем информацию о пользователе
	log.Printf("Пользователь: ID=%d, Имя=%s, Фамилия=%s, Username=%s", client.UserID, client.FirstName, client.LastName, client.Username)

	// Сохраняем данные в JSON-файл
	saveClientData(client)

	// Формируем текст с подставленным именем
	message := fmt.Sprintf(`*Регистрация прошла успешно\! Ваш подарок ждет вас\!*  

Благодарим за регистрацию на вебинар\!  

%s, благодарим за регистрацию на вебинар\!  

Вы сделали важный шаг к повышению эффективности вашей компании через внедрение искусственного интеллекта в CRM\-систему\.  
На вебинаре вы узнаете, как современные инструменты Битрикс24 и AI могут помочь вашему бизнесу\.  

*Дата и время мероприятия:* *27 февраля, 11:00 МСК*  
Ссылка на подключение будет отправлена вам перед мероприятием\.  

*💡 Что вас ждет на вебинаре\?*  

• *Как получить первые быстрые результаты от внедрения CRM\-системы:* цифровизацию, контроль и удобство\.  
• *AI — как реально экономить и зарабатывать:* практические советы по использованию AI для работы в команде, контроля качества, экономии и увеличения продаж\.  
• *Ответы на вопросы:* возможность получить рекомендации от эксперта, бизнес\-аналитика "CRM Консалтинг"\.  

Чтобы помочь вам разобраться в вопросе, мы подготовили для вас *ПОДАРОК*:  

*Гайд "7 возможностей искусственного интеллекта в Битрикс24 для роста продаж"*  

[Забрать подарок](https://drive.google.com/file/d/1gFM1KR9NDqBv2EKLzW_SzWO5ft9qnxhE/view?usp=drive_link)  

До встречи на вебинаре\!  

Если у вас есть вопросы, смело отвечайте на это сообщение — мы всегда рады помочь\.  

С уважением,  
Команда "CRM Консалтинг"`, client.FirstName)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      message,
		ParseMode: "MarkdownV2",
	})
}

// Функция для сохранения данных клиента в JSON-файл (с поддержкой массива)
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

// Функция загрузки клиентов из файла
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
