package commands

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
)

func soy(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	if len(argv) != 1 {
		reply = fmt.Sprintf("Dime quien eres: !soy TuNick")
	} else {
		err := database.SetNick(update.Message.From, argv[0])
		if err != nil {
			if errors.Is(err, database.ErrPandaExists) {
				reply = fmt.Sprintf("Tu ya eres")
			} else {
				reply = fmt.Sprintf("Algo no fue bien: %s", err)
			}
		} else {
			reply = fmt.Sprintf("Vale, a partir de ahora eres <pre>%s</pre>", argv[0])
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func quien(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	if len(argv) != 1 {
		reply = fmt.Sprintf("¿Por quien preguntas?")
	} else {
		username, err := database.TGUserFromNick(argv[0])
		if err != nil {
			reply = fmt.Sprintf("Algo no fue bien: %s", err)
		} else {
			reply = fmt.Sprintf("%s es el panda anteriormente conocido como %s", username, argv[0])
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
