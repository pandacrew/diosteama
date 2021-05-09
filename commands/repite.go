package commands

import (
	"fmt"
	"html"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/format"
)

func repite(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	m := update.Message

	// If its a reply, just use the message
	if m.ReplyToMessage != nil {
		msg := makeMessage(update, []*tgbotapi.Message{m})
		bot.Send(msg)
		return
	}

	// If not, get the forwarded messages with a message queue
	cb := func(q msgQueue) {
		msg := makeMessage(update, q.Messages)
		bot.Send(msg)
	}

	StartMsgQueue(update.Message, cb)
}

func makeMessage(update tgbotapi.Update, msgs []*tgbotapi.Message) tgbotapi.MessageConfig {
	text := format.FormatTGMessages(msgs)
	nick := format.FormatTGUser(update.Message.From)
	formatted := fmt.Sprintf("<pre>%s</pre>\n\n<em>🚽 Quote by %s on %d</em>",
		text, html.EscapeString(nick), format.ParseTime(strconv.Itoa(update.Message.Date)))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, formatted)
	msg.ParseMode = "html"

	return msg
}
