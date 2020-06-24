package commands

import (
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
)

func cmdTop(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var msg tgbotapi.MessageConfig
	var reply string
	var err error

	var i int
	var r string
	if len(argv) == 2 {
		var err error
		i, err = strconv.Atoi(argv[1])
		if err != nil {
			i = 10
		}
	} else {
		i = 10
	}
	r, err = database.Top(i)
	if err != nil {
		log.Println("Error reading top", err)
	}
	reply = strings.Join([]string{"<pre>", r, "</pre>"}, "")
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
