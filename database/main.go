package database

import (
	"log"
	"os"
	"time"

	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool
var addquoteWait time.Duration
var loc *time.Location

func Init() {
	var err error
	dbDsn := os.Getenv("DIOSTEAMA_DB_URL")

	addquoteWait = 800 * time.Millisecond
	loc, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatal(err)
	}

	pool, err = pgxpool.Connect(context.Background(), dbDsn)
	if err != nil {
		log.Panic("Can't create pool", err)
	}

}
