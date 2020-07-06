package format

import (
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
	"github.com/pandacrew-net/diosteama/quotes"
)

func PrettyUser(user *tgbotapi.User) string {
	if user.FirstName != "" {
		var str strings.Builder
		str.WriteString(user.FirstName)
		if user.LastName != "" {
			str.WriteString(" ")
			str.WriteString(user.LastName)
		}
		return str.String()
	}
	return fmt.Sprintf("@%s", user.UserName)
}

func parseTime(t string) time.Time {
	loc, err := time.LoadLocation("Europe/Berlin")
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

// Quote formats a quote to be delivered to the chat
func Quote(quote quotes.Quote) string {
	var nick string
	var err error

	if quote.From == nil {
		nick = strings.SplitN(quote.Author, "!", 2)[0]
	} else {
		nick, err = database.NickFromTGUser(quote.From)
		if err != nil {
			nick = PrettyUser(quote.From)
		}
	}
	//💩🔞🔪💥

	var text string
	log.Printf("%v", quote.Messages)
	if quote.Messages == nil {
		text = html.EscapeString(quote.Text)
	} else {
		text = formatTGMessages(quote.Messages)
	}

	formatted := fmt.Sprintf("<pre>%s</pre>\n\n<em>🚽 Quote %d by %s on %s</em>",
		text, quote.Recnum, html.EscapeString(nick), parseTime(quote.Date))
	return formatted
}

func formatTGMessages(msgs []*tgbotapi.Message) string {
	var result string
	for i := range msgs {
		result = result + formatTGMessage(msgs[i])
	}
	return result
}

func formatTGMessage(msg *tgbotapi.Message) string {
	var user *tgbotapi.User
	var name, text string
	if msg.ReplyToMessage != nil {
		user = msg.ReplyToMessage.From
		text = msg.ReplyToMessage.Text
	} else {
		user = msg.ForwardFrom
		text = msg.Text
	}

	// Uncomment this to use the IRC nick on stored quotes
	name, err := database.NickFromTGUser(user)
	if err != nil {
		name = PrettyUser(user)
	}

	return fmt.Sprintf("%s: %s\n", name, text)
}

// RawQuote creates a string out from a list of raw quotes
func RawQuote(msgs []*tgbotapi.Message) string {
	var result string
	for i := range msgs {
		result = result + RawQuoteMessage(msgs[i])
	}
	return result
}

// RawQuoteMessage creates author: text from a raw message
func RawQuoteMessage(msg *tgbotapi.Message) string {
	var user *tgbotapi.User
	var name, text string
	if msg.ReplyToMessage != nil {
		user = msg.ReplyToMessage.From
		text = msg.ReplyToMessage.Text
	} else {
		user = msg.ForwardFrom
		text = msg.Text
	}

	// Reverted to old behaviour
	name = user.FirstName

	/*
		// Uncomment this to use the IRC nick on stored quotes
		name, err := database.NickFromTGUser(user)
		if err != nil {
			name = user.String()
		}
	*/

	return fmt.Sprintf("%s: %s\n", name, text)
}
