package database

import (
	"fmt"
	"html"
	"log"
	"strings"

	"context"

	"github.com/pandacrew-net/diosteama/format"
	"github.com/pandacrew-net/diosteama/quotes"
)

func nextQuote() int {
	var recnum int
	err := pool.QueryRow(context.Background(), "select max(recnum) from linux_gey_db").Scan(&recnum)
	if err != nil {
		return -1
	}
	recnum = recnum + 1
	if recnum > 20000 {
		return recnum
	}
	return 20000
}

// InsertQuote adds a new quote record to the database
func InsertQuote(quote quotes.Quote) (quotes.Quote, error) {
	var err error

	if quote.Recnum == 0 {
		quote.Recnum = nextQuote()
	}

	query := `
	INSERT INTO linux_gey_db (recnum, date, author, quote, telegram_messages, telegram_author)
	VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb)
	`
	_, err = pool.Exec(context.Background(), query,
		quote.Recnum, quote.Date, quote.Author, quote.Text, quote.Messages, quote.From)

	return quote, err
}

// Info returns all info for a quote
func Info(recnum int, text ...string) (string, error) {
	var quote quotes.Quote
	var query string
	var where string
	var order string
	var parsedQuote string

	query = "SELECT recnum, quote, author, date FROM linux_gey_db"
	where = ""
	if len(text) > 0 {
		where = fmt.Sprintf("WHERE LOWER(quote) LIKE LOWER('%%%s%%')", text[0])
	}

	if recnum < 1 {
		log.Println("Random quote")
		order = "ORDER BY random() LIMIT 1"
	} else {
		where = fmt.Sprintf("WHERE recnum = %d", recnum)
	}
	err := pool.QueryRow(context.Background(),
		fmt.Sprintf("%s %s %s", query, where, order)).Scan(&quote.Recnum, &quote.Text, &quote.Author, &quote.Date)

	if err != nil {
		log.Printf("Error consultando DB: %s", err)
		return "Quote no encontrado", nil
	}
	log.Println(quote.Recnum, quote.Text, quote.Author, quote.Date)

	parsedQuote = format.Quote(quote)
	return parsedQuote, nil
}

// GetQuote performs a quote search
func GetQuote(q string, offset int) (string, error) {
	var b strings.Builder
	var err error
	var count int

	pq := strings.Replace(q, "*", "%", -1)
	query := fmt.Sprintf(`
	SELECT count(*)
	FROM linux_gey_db WHERE LOWER(quote) LIKE LOWER('%%%s%%');`, pq)
	err = pool.QueryRow(context.Background(), query).Scan(&count)
	if err != nil || count < 1 {
		return fmt.Sprintf("Por %s no me sale nada", q), nil
	}

	query = fmt.Sprintf(`
	SELECT recnum, quote
	FROM linux_gey_db WHERE LOWER(quote) LIKE LOWER('%%%s%%')
	ORDER BY recnum ASC LIMIT 5 OFFSET %d;`, pq, offset)
	rows, err := pool.Query(context.Background(), query)
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

// Top performs some unknown stuff (muahaha)
func Top(i int) (string, error) {
	var b strings.Builder
	var err error
	if i < 0 {
		i = 10
	}
	rows, err := pool.Query(context.Background(),
		"select count(*) as c, substring_index(author, '!', 1) as a from linux_gey_db group by a order by c desc limit ?;", i)
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
