package commands

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type commandFunction func(tgbotapi.Update, *tgbotapi.BotAPI, []string)

var invokers = map[string]commandFunction{
	"addquote": addquoteStart,
	"quote":    quote,
	"info":     info,
	"rquote":   rquote,
	"top":      top,
	"culote":   culote,
	"chuches":  chuches,
	"w00g":     w00g,
	"soy":      soy,
	"quienes":  quienes,
	"es":       es,
}

// Command executes a bot command
func Command(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	argv := strings.SplitN(update.Message.Text, " ", 3)
	cmd := argv[0][1:]
	args := argv[1:]

	fn := invokers[cmd]
	fn(update, bot, args)
}
