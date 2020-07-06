package commands

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Command executes a bot command
func Command(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	argv := strings.SplitN(update.Message.Text, " ", 3)
	cmd := argv[0][1:]

	switch cmd {
	case "addquote":
		addquoteStart(update, bot, argv)
	case "quote":
		quote(update, bot, argv)
	case "info":
		info(update, bot, argv)
	case "rquote":
		rquote(update, bot, argv)
	case "top":
		top(update, bot, argv)
	case "culote":
		culote(update, bot, argv)
	case "chuches":
		chuches(update, bot, argv)
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
