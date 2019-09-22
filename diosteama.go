package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var db *sql.DB
var loc *time.Location
func main() {
	var err error
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	dbDsn := os.Getenv("DIOSTEAMA_DB_URL")
	loc, _ = time.LoadLocation("Europe/Andorra")
	db, err = sql.Open("mysql", dbDsn)
	if err != nil {
		log.Panic(err)
	}
	quote("")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		var msg tgbotapi.MessageConfig
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		split := strings.SplitN(update.Message.Text, " ", 2)
		if split[0] == "!quote" || split[0] == "/quote" {
			q := ""
			if len(split) == 2 {
				q = split[1]
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, quote(q))
		} else {
			log.Printf("[%s] %s (%v)", update.Message.From.UserName, update.Message.Text, update.Message.IsCommand())
		}

		/* 		msg.ReplyToMessageID = update.Message.MessageID */
		bot.Send(msg)
	}
}

func quote(q string) string {
	var (
		recnum              int
		date, author, quote string
		f                   string
	)

	query := "SELECT recnum, quote, author, date FROM linux_gey_db"
	if q == "" {
		log.Println("Random quote")
		f = "ORDER BY rand() LIMIT 1"
	} else if i, err := strconv.Atoi(q); err == nil {
		log.Println("Quote by index")
		f = fmt.Sprintf("WHERE recnum = %d", i)
	} else {
		f = fmt.Sprintf("WHERE quote LIKE '%%%s%%' ORDER BY rand() LIMIT 1", q)
	}
	err := db.QueryRow(fmt.Sprintf("%s %s", query, f)).Scan(&recnum, &quote, &author, &date)

	if err != nil {
		return("Quote no encontrado")
	}
	log.Println(recnum, quote, author, date)
	split := strings.SplitN(author, "!", 2)
	return fmt.Sprintf("%s\n\n-- Quote %d by %s on %s", quote, recnum, split[0], parseTime(date))
}

func parseTime(t string) time.Time {
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0).In(loc)
	return tm
}
