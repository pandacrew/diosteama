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
			reply, err := quote(q)
			if err != nil {
				log.Println("Error reading quote: ", err)
			}
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else if split[0] == "!top" || split[0] == "/top" {
			var i int
			if len(split) == 2 {
				var err error
				i, err = strconv.Atoi(split[1])
				if err != nil {
					i = 10
				}
			} else {
				i = 10
			}
			r, err := top(i)
			if err != nil {
				log.Println("Error reading top", err)
			}
			reply := strings.Join([]string{"```", r, "```"}, "")
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ParseMode = "markdown"
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else {
			log.Printf("[%s] %s (%v)", update.Message.From.UserName, update.Message.Text, update.Message.IsCommand())
		}
	}
}

func quote(q string) (string, error) {
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
		return "Quote no encontrado", err
	}
	log.Println(recnum, quote, author, date)
	split := strings.SplitN(author, "!", 2)
	nick := split[0]
	//💩🔞🔪💥
	return fmt.Sprintf("%s\n\n🚽 Quote %d by %s on %s", quote, recnum, nick, parseTime(date)), nil
}

func top(i int) (string, error) {
	var b strings.Builder
	var err error
	if i < 0 {
		i = 10
	}
	rows, err := db.Query("select count(*) as c, substring_index(author, '!', 1) as a from linux_gey_db group by a order by c desc limit ?;", i)
	if err != nil {
		log.Printf("Error listing top %d. Fuck you.", i)
		return b.String(), err
	}
	defer rows.Close()
	i = 0
	for rows.Next() {
		i++
		var (
			count  int
			author string
		)
		err := rows.Scan(&count, &author)
		if err != nil {
			log.Printf("Error scanning top results. Fuck you all!")
			return b.String(), err
		}
		log.Println(count, author)
		fmt.Fprintf(&b, "%3d %20s %5d\n", i, author, count)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("Error in the final possible place in the top 10. Fuck you all! And especially you!")
		return b.String(), err
	}
	return b.String(), err
}

func parseTime(t string) time.Time {
	i, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		i = 1
	}
	tm := time.Unix(i, 0).In(loc)
	return tm
}
