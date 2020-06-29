package commands

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func chuches(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var msg tgbotapi.MessageConfig
	var reply string

	if len(argv) == 1 { // rquote
		reply = fmt.Sprintf("%s, tienes el monopolio de las chuches, no seas avaricioso", update.Message.From.FirstName)
	} else {
		reply = fmt.Sprintf("%s, %s te va a comprar una booolsa de chuuuuches", argv[1], update.Message.From.FirstName)
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	bot.Send(msg)
}
