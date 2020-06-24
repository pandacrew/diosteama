package commands

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pandacrew-net/diosteama/database"
	"github.com/pandacrew-net/diosteama/format"
)

var pool *pgxpool.Pool
var addquotePool map[int]addquote
var addquoteWait time.Duration

type addquote struct {
	UserID   int
	Messages []*tgbotapi.Message
	Timer    *time.Timer
}

func cmdAddquote(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {

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
		addquote := addquote{
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
	addquotePool[uid] = addquote{
		UserID: uid,
		Timer:  time.AfterFunc(addquoteWait, commit),
	}

}

func saveAddquote(uid int, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if existing, exists := addquotePool[uid]; exists {

		var quote database.Quote
		var err error
		var added string
		var msg tgbotapi.MessageConfig

		if len(existing.Messages) < 1 {
			return
		}
		quote.Date = strconv.Itoa(update.Message.Date)
		quote.Text = format.FormatRawQuote(existing.Messages)
		quote.Author = update.Message.From.FirstName // This would be better with a map of telegram users to irc nicks
		quote.Messages = existing.Messages
		quote.From = *update.Message.From

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
