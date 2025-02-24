package webinar

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
	"os/signal"
)

func StartTgBot() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New("780504069:AAH7Ld_hobbvEkCZi8fpdKUIEXirpG4raCQ", opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Регистрация прошла успешно! Ваш подарок ждет вас\nИмя, благодарим за регистрацию за вебинар!\nВы сделали важный шаг к повышению эффективности вашей компании через внедрение искусственного интеллекта в CRM-систему. На вебинаре вы узнаете, как современные инструменты Битрикс24 и AI могут помочь вашему бизнесу.\nДата и время мероприятия: 27 февраля, 11:00 МСК. Ссылка на подключение будет отправлена вам перед мероприятием.\n💡 Что вас ждет на вебинаре?\nКак получить первые быстрые результат от внедрения CRM-системы: цифровизацию, контроль и удобство. \nAI — как реально экономить и зарабатывать: практические советы по использованию AI для работы в команде, контроля качества, экономии и увеличения продаж.\nОтветы на вопросы: возможность получить рекомендации от эксперта, бизнес-аналитика “CRM Консалтинг”.\nЧтобы помочь вам разобраться в вопросе, мы подготовили для вас ПОДАРОК:\nГайд “7 возможностей искусственного интеллекта в Битрикс24 для роста продаж”.\n Забрать подарок https://t.me/crmconsaltigbot \nДо встречи на вебинаре!\nЕсли у вас есть вопросы, смело отвечайте на это сообщение — мы всегда рады помочь.\nС уважением,\nКоманда “CRM Консалтинг”\n",
	})
}
