package commands

import (
	"fmt"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func chuches(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string

	if len(argv) == 0 { // rquote
		reply = fmt.Sprintf("%s, tienes el monopolio de las chuches, no seas avaricioso", update.Message.From.FirstName)
	} else {
		reply = fmt.Sprintf("%s, %s te va a comprar una booolsa de chuuuuches",
			strings.Join(argv, " "), update.Message.From.FirstName)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	bot.Send(msg)
}
