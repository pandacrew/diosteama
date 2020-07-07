package quotes

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Quote is the full quote
type Quote struct {
	Recnum   int
	Date     string
	Author   string
	Text     string
	Messages []*tgbotapi.Message
	From     *tgbotapi.User
}
