package main

import (
	"github.com/pandacrew-net/diosteama/database"
	"github.com/pandacrew-net/diosteama/telegram"
)

func main() {
	database.Init()
	database.Info(0)
	telegram.Start()
}
