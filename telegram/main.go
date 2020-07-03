package telegram

import (
	"log"
	"os"

	"encoding/json"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/commands"
)

// Start initialized the bot and runs main loop
func Start() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		j, _ := json.Marshal(update)
		log.Printf("%s", j)
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		response(update, bot)
		log.Printf("[%s] %s (%v)", update.Message.From.UserName, update.Message.Text, update.Message.IsCommand())

	}
}

func response(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	if commands.EvalAddquote(update) {
		// This is a forward part of an !addquote and has been processed. Return.
		return
	}

	if len(update.Message.Text) > 0 && (string(update.Message.Text[0]) == "!" || string(update.Message.Text[0]) == "/") {
		commands.Command(update, bot)
		return
	}

	commands.Triggers(update, bot)
}
