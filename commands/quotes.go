package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pandacrew-net/diosteama/database"
	"github.com/pandacrew-net/diosteama/format"
	"github.com/pandacrew-net/diosteama/quotes"
)

var pool *pgxpool.Pool

func addquote(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	m := update.Message

	// If its a reply, just use the message
	if m.ReplyToMessage != nil {
		q := msgQueue{
			UserID:   m.From.ID,
			Messages: []*tgbotapi.Message{m},
		}
		saveQuotes(update, bot, q)
		return
	}

	// If not, get the forwarded messages with a message queue
	cb := func(q msgQueue) {
		saveQuotes(update, bot, q)
	}

	StartMsgQueue(update.Message, cb)
}

func saveQuotes(update tgbotapi.Update, bot *tgbotapi.BotAPI, q msgQueue) {
	var quote quotes.Quote
	var err error
	var added string
	var msg tgbotapi.MessageConfig

	if len(q.Messages) < 1 {
		log.Printf("No messages to save")
		return
	}
	quote.Date = strconv.Itoa(update.Message.Date)
	quote.Text = format.RawQuote(q.Messages)
	quote.Author = format.PrettyUser(update.Message.From) // This would be better with a map of telegram users to irc nicks
	quote.Messages = q.Messages
	quote.From = update.Message.From

	quote, err = database.InsertQuote(quote)

	if err != nil {
		//time.Sleep(addquoteWait)
		//saveAddquote(uid, update, bot)
		log.Fatalf("Error saving quote %d: %v", quote.Recnum, err)
	}

	log.Printf("Saved quote %d for %d, %s, %d", quote.Recnum, q.UserID, quote.From, update.Message.Date)

	added = fmt.Sprintf("Quote added: %d", quote.Recnum)
	log.Println(added)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, added)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func quote(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	var err error

	if len(argv) == 0 { // rquote
		var quote *quotes.Quote
		quote, err = database.Info(-1)
		if err != nil {
			log.Println("Error reading quote: ", err)
		} else {
			reply = format.Quote(*quote)
		}
	} else {
		offset, err := strconv.Atoi(argv[0])
		if err != nil || len(argv) == 1 || offset < 0 {
			text := strings.Join(argv, " ")
			reply, err = database.GetQuote(text, 0)
		} else {
			reply, err = database.GetQuote(argv[1], offset)
		}
		if err != nil {
			log.Println("Error reading quote: ", err)
		}
	}
	log.Println("Replying", reply)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func rquote(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var quote *quotes.Quote
	var reply string
	var err error

	if len(argv) == 0 {
		quote, err = database.Info(-1)
	} else if len(argv) == 1 {
		quote, err = database.Info(-1, argv[0])
	}
	if err != nil {
		log.Println("Error reading quote: ", err)
	}
	reply = format.Quote(*quote)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func info(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string

	if len(argv) < 1 {
		reply = "Error. Format is !info <quote id>"
	}

	qid, err := strconv.Atoi(argv[0])
	if err != nil {
		reply = "Error. Format is !info <quote id>"
	}

	quote, err := database.Info(qid)
	if err != nil {
		log.Println("Error reading quote: ", err)
	} else {
		reply = format.Quote(*quote)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
