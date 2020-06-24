package commands

import (
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Command executes a bot command
func Command(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	argv := strings.SplitN(update.Message.Text, " ", 3)
	cmd := argv[0][1:]

	switch cmd {
	case "addquote":
		cmdAddquote(update, bot, argv)
	case "quote":
		cmdQuote(update, bot, argv)
	case "info":
		cmdInfo(update, bot, argv)
	case "rquote":
		cmdRquote(update, bot, argv)
	case "top":
		cmdTop(update, bot, argv)
	case "culote":
		cmdCulote(update, bot, argv)
	case "chuches":
		cmdChuches(update, bot, argv)
	case "w00g":
		cmdW00g(update, bot, argv)
	}
}
