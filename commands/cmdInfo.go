package commands

import (
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
)

func cmdInfo(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var msg tgbotapi.MessageConfig
	var reply string
	var quote string
	var err error

	if len(argv) < 2 {
		reply = "Error. Format is !info <quote id>"
	}

	qid, err := strconv.Atoi(argv[1])
	if err != nil {
		reply = "Error. Format is !info <quote id>"
	}

	quote, err = database.Info(qid)
	if err != nil {
		log.Println("Error reading quote: ", err)
	} else {
		reply = quote
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
