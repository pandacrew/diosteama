package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn
var loc *time.Location

func command(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg tgbotapi.MessageConfig
	var reply string
	split := strings.SplitN(update.Message.Text, " ", 3)
	switch cmd := split[0][1:]; cmd {
	case "quote":
		var err error
		if len(split) == 1 { // rquote
			reply, err = info(-1)
			if err != nil {
				log.Println("Error reading quote: ", err)
			}
		} else if len(split) == 2 {
			reply, err = quote(split[1], 0)
			if err != nil {
				log.Println("Error reading quote: ", err)
			}
		} else {
			offset, err := strconv.Atoi(split[1])
			if err != nil || offset < 0 {
				reply = "Error. Format is <code>!quote [[offset] search]</code>"
			} else {
				reply, err = quote(split[2], offset)
				if err != nil {
					log.Println("Error reading quote: ", err)
				}
			}
		}
		log.Println("Replying", reply)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "html"
		bot.Send(msg)
	case "info":
		if len(split) < 2 {
			reply = "Error. Format is !info <quote id>"
		}
		qid, err := strconv.Atoi(split[1])
		if err != nil {
			reply = "Error. Format is !info <quote id>"

		}
		reply, err = info(qid)
		if err != nil {
			log.Println("Error reading quote: ", err)

		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "html"
		bot.Send(msg)
	case "rquote":
		reply, err := info(-1)
		if err != nil {
			log.Println("Error reading quote: ", err)

		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "html"
		bot.Send(msg)
	case "top":
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
		reply = strings.Join([]string{"<pre>", r, "</pre>"}, "")
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "html"
		bot.Send(msg)
	case "culote":
		reply = fmt.Sprintf("%s, tienes un culote como para meter %s", update.Message.From.FirstName, strings.Join(split[1:], " "))
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	case "chuches":
		if len(split) == 1 { // rquote
			reply = fmt.Sprintf("%s, tienes el monopolio de las chuches, no seas avaricioso", update.Message.From.FirstName)
		} else {
			reply = fmt.Sprintf("%s, %s te va a comprar una booolsa de chuuuuches", split[1], update.Message.From.FirstName)
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}

}

func main() {
	var err error
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	dbDsn := os.Getenv("DIOSTEAMA_DB_URL")
	loc, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatal(err)
	}

	conn, err = pgx.Connect(context.Background(), dbDsn)
	if err != nil {
		log.Panic(err)
	}
	info(0)
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
		j, _ := json.Marshal(update)
		log.Printf("%s", j)
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		response(update, bot)
		log.Printf("[%s] %s (%v)", update.Message.From.UserName, update.Message.Text, update.Message.IsCommand())

	}
}

func eval_addquote(msg tgbotapi.Update) {
}

func response(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var msg tgbotapi.MessageConfig
	if string(update.Message.Text[0]) == "!" || string(update.Message.Text[0]) == "/" {
		command(update, bot)
	} else if strings.Contains(strings.ToLower(update.Message.Text), "almeida") {
		reply := "¡¡CARAPOLLA!!"
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	} else if strings.Contains(strings.ToLower(update.Message.Text), "carme") {
		reply := "PUTAAAAAAAAAA"
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	} else if strings.Contains(strings.ToLower(update.Message.Text), "gamba") {
		reply := "MARIPURI!"
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
func info(i int) (string, error) {
	var (
		recnum              int
		date, author, quote string
		f                   string
	)

	query := "SELECT recnum, quote, author, date FROM linux_gey_db"
	if i < 1 {
		log.Println("Random quote")
		f = "ORDER BY random() LIMIT 1"
	} else {
		f = fmt.Sprintf("WHERE recnum = %d", i)
	}
	err := conn.QueryRow(context.Background(), fmt.Sprintf("%s %s", query, f)).Scan(&recnum, &quote, &author, &date)

	if err != nil {
		log.Printf("Error consultando DB: %s", err)
		return "Quote no encontrado", nil
	}
	log.Println(recnum, quote, author, date)
	split := strings.SplitN(author, "!", 2)
	nick := split[0]
	//💩🔞🔪💥
	return fmt.Sprintf("<pre>%s</pre>\n\n<em>🚽 Quote %d by %s on %s</em>", html.EscapeString(quote), recnum, html.EscapeString(nick), parseTime(date)), nil
}

func quote(q string, offset int) (string, error) {
	var b strings.Builder
	var err error
	var count int
	pq := strings.Replace(q, "*", "%", -1)
	query := fmt.Sprintf(`
	SELECT count(*)
	FROM linux_gey_db WHERE LOWER(quote) LIKE LOWER('%%%s%%');`, pq)
	err = conn.QueryRow(context.Background(), query).Scan(&count)
	if err != nil || count < 1 {
		return fmt.Sprintf("Por %s no me sale nada", q), nil
	}

	query = fmt.Sprintf(`
	SELECT recnum, quote
	FROM linux_gey_db WHERE LOWER(quote) LIKE LOWER('%%%s%%')
	ORDER BY recnum ASC LIMIT 5 OFFSET %d;`, pq, offset)
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error getting quotes for %s. Fuck you.", q)
		return b.String(), err
	}
	defer rows.Close()
	i := offset

	for rows.Next() {
		i++
		var (
			recnum int
			quote  string
		)
		err := rows.Scan(&recnum, &quote)
		if err != nil {
			log.Printf("Error getting quotes. Fuck you all!")
			return b.String(), err
		}
		fmt.Fprintf(&b, "%d. <code>%s</code>\n", recnum, html.EscapeString(quote))
	}
	fmt.Fprintf(&b, "\nQuotes %d a %d de %d buscando <code>%s</code>", offset+1, i, count, html.EscapeString(q))
	err = rows.Err()
	if err != nil {
		log.Printf("Error in the final possible place getting quotes. Fuck you all! And especially you!")
		return b.String(), err
	}
	log.Println(b.String())
	return b.String(), err
}

func top(i int) (string, error) {
	var b strings.Builder
	var err error
	if i < 0 {
		i = 10
	}
	rows, err := conn.Query(context.Background(), "select count(*) as c, substring_index(author, '!', 1) as a from linux_gey_db group by a order by c desc limit ?;", i)
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
