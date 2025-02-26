package webinar

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
)

/*
func SendMessageIn9AM(b *bot.Bot) {
	clientName := "Danny"
	message := fmt.Sprintf(`*Вебинар через 2 часа\! Как CRM \+ AI увеличивают прибыль\?*

%s, доброе утро\!

На связи команда CRM Consulting, через 2 часа мы разберём, как автоматизировать продажи, перестать терять клиентов и увеличить прибыль с помощью CRM и искусственного интеллекта\.

Вот что вы узнаете на вебинаре:
• Как CRM \+ AI помогают не терять до 40%% клиентов – автоматизация заявок, напоминания, персонализированные follow\-up\'ы\.
• Почему CRM – это не просто база данных, а инструмент роста – реальные кейсы, как бизнесы увеличили выручку на 30\-50%%\.
• Как построить отдел продаж, который работает без хаоса – CRM поможет внедрить чёткие процессы\.

Подключайтесь вовремя, чтобы не пропустить самое важное\!
👉 [Присоединяйтесь к вебинару здесь](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)

Скоро увидимся в эфире\!
Команда CRM Consulting`, clientName)

	ctx := context.Background()
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    testUserID,
		Text:      message,
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("Ошибка при отправке тестового сообщения: %v", err)
	} else {
		log.Println("Тестовое сообщение успешно отправлено пользователю")
	}
}

func SendMessageIn1050AM(b *bot.Bot) {
	clientName := "Danny"
	message := fmt.Sprintf(`*Осталось 10 минут и начинаем\. Подключайтесь*

%s, уже через 10 минут разберем, как CRM \+ AI могут закрыть ключевые проблемы бизнеса и помочь вам зарабатывать больше\!

Вам будет полезно, если\:

• *Менеджеры теряют клиентов* – заявки обрабатываются с опозданием, клиенты уходят к конкурентам\.
• *Нет прозрачности в продажах* – сложно понять, на каком этапе сделки, кто выполняет задачи и где узкие места\.
• *Слабый контроль за командой* – не хватает данных, чтобы объективно оценить работу менеджеров\.
• *Ручная работа замедляет бизнес* – отчёты, напоминания, письма менеджеры делают вручную вместо продаж\.
• CRM есть, но работает вполсилы или вы задумываетесь о внедрении\.

🔥 Сегодня разберём, как CRM \+ AI решают эти проблемы и превращают хаос в систему, которая реально продает\!

[Ваша ссылка на подключение](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)`, clientName)

	ctx := context.Background()
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    testUserID,
		Text:      message,
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("Ошибка при отправке тестового сообщения: %v", err)
	} else {
		log.Println("Тестовое сообщение успешно отправлено пользователю")
	}
}

func SendMessageIn11AM(b *bot.Bot) {

	message := fmt.Sprintf(`*3\.\.\.2\.\.\.1\.\.\. Мы в эфире\! Подключайтесь прямо сейчас*

Мы уже начали вебинар\! *Как CRM \+ AI увеличивают продажи, сокращают потери клиентов и систематизируют бизнес\?* Разберем прямо сейчас\!

Обсудим без воды: как перестать терять клиентов и увеличить прибыль с CRM, какие процессы можно автоматизировать уже сейчас, реальные кейсы компаний,
которые уже используют наши клиенты\.

Присоединяйтесь: 👉 [Ссылка на вебинар](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)

Не упустите, то что для вас важно – присоединяйтесь\!`)

	ctx := context.Background()
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    testUserID,
		Text:      message,
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("Ошибка при отправке тестового сообщения: %v", err)
	} else {
		log.Println("Тестовое сообщение успешно отправлено пользователю")
	}
}

func SendMessageIn1115AM(b *bot.Bot) {

	message := fmt.Sprintf(`*Как подружить CRM и искусственный интеллект\?*

Михаил, бизнес\-аналитик CRM Consulting рассказывает в эфире прямо сейчас, сэкономить и увеличить продажи в бизнесе с помощью искусственного интеллекта\.

Заходите: [Ссылка](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)
`)

	ctx := context.Background()
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    testUserID,
		Text:      message,
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("Ошибка при отправке тестового сообщения: %v", err)
	} else {
		log.Println("Тестовое сообщение успешно отправлено пользователю")
	}
}

func SendMessageIn1130AM(b *bot.Bot) {

	message := fmt.Sprintf(`*🔥 Снизили затраты на Call\-центр в 10 раз\! С помощью ИИ*

Сегодня на вебинаре мы рассматриваем реальные истории внедрения ИИ \+ CRM\.

Подключайтесь, чтобы узнать, как компания отказалась от Call\-центра на аутсорсе благодаря искусственному интеллекту\.

Ждем вас: 👉 [Ссылка на вебинар](https://us06web.zoom.us/j/89326719994?pwd=GysR7TlRgW0l1lSIUqGEzj016yZ1jU.1)
`)

	ctx := context.Background()
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    testUserID,
		Text:      message,
		ParseMode: "MarkdownV2",
	})

	if err != nil {
		log.Printf("Ошибка при отправке тестового сообщения: %v", err)
	} else {
		log.Println("Тестовое сообщение успешно отправлено пользователю")
	}
}*/

func Handler(ctx context.Context, b *bot.Bot, update *models.Update) {
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
