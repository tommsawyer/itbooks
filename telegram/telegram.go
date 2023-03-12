package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var api *tgbotapi.BotAPI

// Authorize will authorize in telegram using provided bot token.
//
// You can obtain new token using telegram bot father.
// More info here: https://core.telegram.org/bots/tutorial#obtain-your-bot-token
func Authorize(ctx context.Context, token string) error {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	api = botAPI
	return nil
}

// Send sends message to provided channel.
// It is required to add bot to that channel.
func Send(ctx context.Context, channel string, msg Message) error {
	imageWithCaption, err := msg.imageWithCaption(channel)
	if err != nil {
		return fmt.Errorf("cannot create image with caption: %w", err)
	}
	_, err = api.Send(imageWithCaption)
	return err
}
