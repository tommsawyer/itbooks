package telegram

import (
	"strings"
	"testing"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/matryer/is"
)

func TestMessageBuildedCorrectly(t *testing.T) {
	is := is.New(t)

	msg := Message{
		ImageURL: "imageurl",
		Title:    "title",
		Subtitle: "subtitle",
		Link:     "link",
		Text:     "text",
	}

	upload, err := msg.imageWithCaption("channel")
	is.NoErr(err)

	is.Equal(upload.BaseChat.ChannelUsername, "channel")
	is.Equal(upload.BaseFile.File.SendData(), "imageurl")
	is.Equal(upload.ParseMode, tgbotapi.ModeMarkdownV2)

	expectedText := `
*title*
_subtitle_
[Купить](link)

text
`

	is.Equal(strings.TrimSpace(upload.Caption), strings.TrimSpace(expectedText))
}

func TestMessageTruncateTooLongTexts(t *testing.T) {
	is := is.New(t)

	msg := Message{
		ImageURL: "imageurl",
		Title:    "title",
		Subtitle: "subtitle",
		Link:     "link",
		Text:     strings.Repeat("test", 1000),
	}

	upload, err := msg.imageWithCaption("channel")
	is.NoErr(err)

	is.Equal(utf8.RuneCountInString(upload.Caption), maxTelegramMessageSize)
}
