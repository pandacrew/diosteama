package commands

import (
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
)

func top(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	i := 10

	if len(argv) == 1 {
		if j, err := strconv.Atoi(argv[0]); err == nil {
			i = j
		}
	}
	r, err := database.Top(i)
	if err != nil {
		log.Println("Error reading top", err)
	}
	reply := strings.Join([]string{"<pre>", r, "</pre>"}, "")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
