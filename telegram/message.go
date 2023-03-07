package telegram

import (
	"embed"
	"strings"
	"text/template"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func (msg Message) photoUpload(channel string) (tgbotapi.Chattable, error) {
	text, err := msg.markdown()
	if err != nil {
		return nil, err
	}

	upload := tgbotapi.NewPhotoToChannel(
		channel,
		tgbotapi.FileURL(msg.ImageURL),
	)
	upload.Caption = text
	upload.ParseMode = tgbotapi.ModeMarkdownV2

	return upload, nil
}

func (msg Message) markdown() (string, error) {
	var builder strings.Builder

	if err := messageTemplates.ExecuteTemplate(&builder, "book.md", msg); err != nil {
		return "", err
	}

	return msg.truncate(builder.String()), nil
}

func (msg Message) truncate(s string) string {
	const maxTelegramMessageSize = 1024

	if utf8.RuneCountInString(s) <= maxTelegramMessageSize {
		return s
	}

	return string([]rune(s)[:maxTelegramMessageSize-3]) + tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, "...")
}
