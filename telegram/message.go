package telegram

import (
	"embed"
	"strings"
	"text/template"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxTelegramMessageSize = 1024

var (
	//go:embed templates/*
	templatesDir embed.FS

	messageTemplates = template.Must(
		template.New("message").
			Funcs(template.FuncMap{
				"escape": func(s string) string {
					return tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, s)
				},
			}).
			ParseFS(templatesDir, "templates/*.md"),
	)
)

// Message is message to telegram.
type Message struct {
	ImageURL string
	Title    string
	Subtitle string
	Link     string
	Text     string
}

func (msg *Message) imageWithCaption(channel string) (tgbotapi.PhotoConfig, error) {
	text, err := msg.markdown()
	if err != nil {
		return tgbotapi.PhotoConfig{}, err
	}

	upload := tgbotapi.NewPhotoToChannel(
		channel,
		tgbotapi.FileURL(msg.ImageURL),
	)
	upload.Caption = text
	upload.ParseMode = tgbotapi.ModeMarkdownV2

	return upload, nil
}

func (msg *Message) markdown() (string, error) {
	var builder strings.Builder

	if err := messageTemplates.ExecuteTemplate(&builder, "book.md", msg); err != nil {
		return "", err
	}

	return msg.truncate(builder.String()), nil
}

func (msg *Message) truncate(s string) string {
	if utf8.RuneCountInString(s) <= maxTelegramMessageSize {
		return s
	}

	dots := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, "...")

	return string([]rune(s)[:maxTelegramMessageSize-utf8.RuneCountInString(dots)]) + dots
}
