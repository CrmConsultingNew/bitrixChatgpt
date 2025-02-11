package webinar

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"net/http"
	"os"
	"os/signal"
)

func StartTgBot() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handlerTgBot),
		bot.WithWebhookSecretToken(os.Getenv("780504069:AAH7Ld_hobbvEkCZi8fpdKUIEXirpG4raCQ")),
	}

	b, _ := bot.New(os.Getenv("780504069:AAH7Ld_hobbvEkCZi8fpdKUIEXirpG4raCQ"), opts...)

	// call methods.SetWebhook if needed

	go b.StartWebhook(ctx)
	http.ListenAndServe(":2000", b.WebhookHandler())
}

func handlerTgBot(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
