package commands

import (
	"log"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
)

func cmdQuote(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var msg tgbotapi.MessageConfig
	var reply string
	var err error
	var offset int
	if len(argv) == 1 { // rquote
		reply, err = database.Info(-1)
		if err != nil {
			log.Println("Error reading quote: ", err)
		}
	} else if len(argv) == 2 {
		reply, err = database.GetQuote(argv[1], 0)
		if err != nil {
			log.Println("Error reading quote: ", err)
		}
	} else {
		offset, err = strconv.Atoi(argv[1])
		if err != nil || offset < 0 {
			reply = "Error. Format is <code>!quote [[offset] search]</code>"
		} else {
			reply, err = database.GetQuote(argv[2], offset)
			if err != nil {
				log.Println("Error reading quote: ", err)
			}
		}
	}
	log.Println("Replying", reply)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
