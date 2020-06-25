package format

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/quotes"
)

func parseTime(t string) time.Time {
	var loc *time.Location
	var err error
	loc, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Println(err)
	}
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		i = 1
	}
	tm := time.Unix(i, 0).In(loc)
	return tm
}

func FormatQuote(quote quotes.Quote) string {
	var nick string
	var formatted string

	nick = strings.SplitN(quote.Author, "!", 2)[0]
	//💩🔞🔪💥
	formatted = fmt.Sprintf("<pre>%s</pre>\n\n<em>🚽 Quote %d by %s on %s</em>",
		html.EscapeString(quote.Text), quote.Recnum, html.EscapeString(nick), parseTime(quote.Date))
	return formatted
}

func FormatRawQuote(msgs []*tgbotapi.Message) string {
	var result string
	for i := range msgs {
		result = result + FormatRawQuoteMessage(msgs[i])
	}
	return result
}

func FormatRawQuoteMessage(msg *tgbotapi.Message) string {
	var name, text string
	if msg.ReplyToMessage != nil {
		name = msg.ReplyToMessage.From.FirstName
		text = msg.ReplyToMessage.Text
	} else {
		name = msg.ForwardFrom.FirstName
		text = msg.Text
	}
	return fmt.Sprintf("%s: %s\n", name, text)
}
