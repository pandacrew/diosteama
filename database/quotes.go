package database

import (
	"errors"
	"fmt"
	"html"
	"log"
	"strings"

	"context"

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
	if quote.Recnum == 0 {
		quote.Recnum = nextQuote()
	}

	query := `
	INSERT INTO linux_gey_db (recnum, date, author, quote, telegram_messages, telegram_author)
	VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb)
	`
	_, err := pool.Exec(context.Background(), query,
		quote.Recnum, quote.Date, quote.Author, quote.Text, quote.Messages, quote.From)

	return quote, err
}

// Info returns all info for a quote
func Info(recnum int, text ...string) (*quotes.Quote, error) {
	var quote quotes.Quote
	var order string

	query := `SELECT recnum, quote, author, date, telegram_messages, telegram_author
		 FROM linux_gey_db WHERE deleted is null`
	where := ""
	if len(text) > 0 {
		where = fmt.Sprintf("AND LOWER(quote) LIKE LOWER('%%%s%%')", text[0])
	}

	if recnum < 1 {
		log.Println("Random quote")
		order = "ORDER BY random() LIMIT 1"
	} else {
		where = fmt.Sprintf("AND recnum = %d", recnum)
	}
	err := pool.QueryRow(context.Background(),
		fmt.Sprintf("%s %s %s", query, where, order)).Scan(&quote.Recnum, &quote.Text, &quote.Author, &quote.Date, &quote.Messages, &quote.From)

	if err != nil {
		log.Printf("Error consultando DB: %s", err)
		return nil, fmt.Errorf("%w", err)
	}
	log.Println(quote.Recnum, quote.Text, quote.Author, quote.Date)

	return &quote, nil
}

// GetQuote performs a quote search
func GetQuote(q string, offset int) (string, error) {
	var b strings.Builder
	var count int

	pq := strings.Replace(q, "*", "%", -1)
	query := fmt.Sprintf(`
	SELECT count(*)
	FROM linux_gey_db WHERE deleted is null AND LOWER(quote) LIKE LOWER('%%%s%%');`, pq)
	err := pool.QueryRow(context.Background(), query).Scan(&count)
	if err != nil || count < 1 {
		return fmt.Sprintf("Por %s no me sale nada", q), nil
	}

	query = fmt.Sprintf(`
	SELECT recnum, quote
	FROM linux_gey_db WHERE deleted is null AND LOWER(quote) LIKE LOWER('%%%s%%')
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
	if i < 0 {
		i = 10
	}
	query := "select count(*) as c, split_part(author, '!', 1) as a from linux_gey_db where deleted is null group by a order by c desc limit $1;"
	rows, err := pool.Query(context.Background(), query, i)
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

// MarkQuoteAsDeleted marks a quote with identifier id as deleted
func MarkQuoteAsDeleted(recnum int, user string) error {
	quote, err := FindQuoteById(recnum, false)
	if err != nil {
		log.Printf("[MarkQuoteAsDeleted] Quote with id %d wasn't found", recnum)
		return err
	}
	const updateStmt = `UPDATE linux_gey_db 
		SET deleted = current_timestamp, deleted_by = $2
		WHERE deleted is null AND recnum = $1;`

	_, err = pool.Exec(context.Background(), updateStmt, quote.Recnum, user)

	return err
}

// UnmarkQuoteAsDeleted marks a deleted quote with identifier id as undeleted
func UnmarkQuoteAsDeleted(recnum int) error {
	quote, err := FindQuoteById(recnum, true)
	if err != nil {
		log.Printf("[UnmarkQuoteAsDeleted] Quote with id %d wasn't found", recnum)
		return err
	}
	const updateStmt = `UPDATE linux_gey_db 
		SET deleted = null, deleted_by = null
		WHERE deleted is not null AND recnum = $1;`

	_, err = pool.Exec(context.Background(), updateStmt, quote.Recnum)

	return err
}

var errDontMess = errors.New("don't mess with me! AKA no me toques lo que no suena")

// FindQuoteById Finds a quote by it's unique id
// recnum: quote id to find
// includeDeleted: true to inlcude deleted quotes in search
func FindQuoteById(recnum int, includeDeleted bool) (*quotes.Quote, error) {
	if recnum < 1 {
		return nil, errDontMess
	}

	var findQuoteByIdQuery = `
		SELECT recnum, quote, author, date, telegram_messages, telegram_author
		FROM linux_gey_db 
		WHERE recnum = $1`

	if !includeDeleted {
		findQuoteByIdQuery += " AND deleted is null"
	}

	var quote quotes.Quote
	err := pool.QueryRow(context.Background(), findQuoteByIdQuery, recnum).
		Scan(&quote.Recnum, &quote.Text, &quote.Author, &quote.Date, &quote.Messages, &quote.From)

	if err != nil {
		log.Printf("[FindQuoteById] Error consultando DB: %v", err)
		return nil, fmt.Errorf("%w", err)
	}
	log.Println(quote.Recnum, quote.Text, quote.Author, quote.Date)

	return &quote, nil
}
