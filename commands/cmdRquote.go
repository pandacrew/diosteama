package commands

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
)

func cmdRquote(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var msg tgbotapi.MessageConfig
	var reply string
	var err error

	if len(argv) == 1 {
		reply, err = database.Info(-1)
	} else if len(argv) == 2 {
		reply, err = database.Info(-1, argv[1])
	}
	if err != nil {
		log.Println("Error reading quote: ", err)
	}
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
