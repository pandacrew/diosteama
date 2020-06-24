package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func cmdW00g(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var msg tgbotapi.MessageConfig
	var reply string

	reply = "Capitan castor, ayuditaaaaaaaaaaaaaaaa!!!"
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	bot.Send(msg)
}
