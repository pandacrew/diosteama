package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pandacrew-net/diosteama/database"
	"github.com/pandacrew-net/diosteama/format"
	"github.com/pandacrew-net/diosteama/quotes"
)

var pool *pgxpool.Pool
var addquotePool map[int]Addquote
var addquoteWait = 800 * time.Millisecond

type Addquote struct {
	UserID   int
	Messages []*tgbotapi.Message
	Timer    *time.Timer
}

func addquoteStart(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	addquotePool = make(map[int]Addquote)
	uid := update.Message.From.ID
	if update.Message.ForwardDate > 0 {
		return
	}
	if existing, exists := addquotePool[uid]; exists {
		// Stop timer for previous addquote, save and start a new one
		existing.Timer.Stop()
		saveAddquote(uid, update, bot)
	}
	if update.Message.ReplyToMessage != nil {
		addquote := Addquote{
			UserID: uid,
		}
		addquote.Messages = append(addquotePool[uid].Messages, update.Message)
		addquotePool[uid] = addquote
		saveAddquote(uid, update, bot)
		return
	}
	commit := func() {
		log.Printf("Expired timer for %d, %s, %s", uid, update.Message.From, update.Message.Date)
		saveAddquote(uid, update, bot)
	}
	addquotePool[uid] = Addquote{
		UserID: uid,
		Timer:  time.AfterFunc(addquoteWait, commit),
	}

}

func saveAddquote(uid int, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if existing, exists := addquotePool[uid]; exists {

		var quote quotes.Quote
		var err error
		var added string
		var msg tgbotapi.MessageConfig

		if len(existing.Messages) < 1 {
			log.Printf("No messages to save")
			return
		}
		quote.Date = strconv.Itoa(update.Message.Date)
		quote.Text = format.RawQuote(existing.Messages)
		quote.Author = update.Message.From.FirstName // This would be better with a map of telegram users to irc nicks
		quote.Messages = existing.Messages
		quote.From = update.Message.From

		quote, err = database.InsertQuote(quote)

		if err != nil {
			//time.Sleep(addquoteWait)
			//saveAddquote(uid, update, bot)
			log.Fatalf("Error saving quote %d: %v", quote.Recnum, err)
		}

		log.Printf("Saved quote %d for %d, %s, %s", quote.Recnum, uid, quote.From, update.Message.Date)

		added = fmt.Sprintf("Quote added: %d", quote.Recnum)
		log.Println(added)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, added)
		msg.ParseMode = "html"
		bot.Send(msg)
		delete(addquotePool, uid)
		log.Printf("Cleanup of addquotePool[%d]", uid)
	} else {
		log.Printf("weird error condition, we were called without an existing pool")
	}
}

// EvalAddquote checks if message is a forward part of an addquote and has been processed
func EvalAddquote(update tgbotapi.Update) bool {
	uid := update.Message.From.ID
	if existing, exists := addquotePool[uid]; exists && update.Message.ForwardDate > 0 {
		existing.Timer.Reset(addquoteWait)
		existing.Messages = append(existing.Messages, update.Message)
		addquotePool[uid] = existing
		return true
	}
	return false
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
		reply = fmt.Sprintf("Quote %d not found", qid);
	} else {
		reply = format.Quote(*quote)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func removeQuote(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	if len(argv) < 1 {
		reply = "Error. Format is !rmquote <quote id>"
	}

	quoteId, err := strconv.Atoi(argv[0])
	if err != nil {
		reply = "Error. Format is !rmquote <quote id>"
	}

	err = database.MarkQuoteAsRemoved(quoteId)
	if err != nil {
		log.Printf("Error removing quote %d: %v", quoteId, err)
		reply = "Error. Quote wasn't removed due errors"
	} else {
		reply = fmt.Sprintf("Quote %d removed!", quoteId)
	}


	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
