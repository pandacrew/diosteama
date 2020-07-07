package commands

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Command executes a bot command
func Command(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	argv := strings.SplitN(update.Message.Text, " ", 3)
	cmd := argv[0][1:]
	args := argv[1:]

	switch cmd {
	case "addquote":
		addquoteStart(update, bot, args)
	case "quote":
		quote(update, bot, args)
	case "info":
		info(update, bot, args)
	case "rquote":
		rquote(update, bot, args)
	case "top":
		top(update, bot, args)
	case "culote":
		culote(update, bot, args)
	case "chuches":
		chuches(update, bot, args)
	case "w00g":
		w00g(update, bot, argv)
	case "soy":
		soy(update, bot, argv)
	case "quienes":
		quienes(update, bot, argv)
	case "es":
		es(update, bot, argv)

	}
}
