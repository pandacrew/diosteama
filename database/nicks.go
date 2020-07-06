package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ErrPandaNotFound error returned when a user or nick is not in the DB
var ErrPandaNotFound = errors.New("No such Panda")

// NoneID null id, sort of.
const NoneID = -1

// NickFromTGUser return the nick for a given Telegram User
func NickFromTGUser(user *tgbotapi.User) (string, error) {
	var nick string

	query := fmt.Sprintf(`SELECT nick FROM user_nick WHERE tg_user_id = $1;`)
	err := pool.QueryRow(context.Background(), query, user.ID).Scan(nick)
	if err != nil {
		log.Printf("Error finding panda %s. Fuck you.", user.FirstName)
		return "", err
	}

	if nick == "" {
		log.Printf("%s is not a panda.", user.FirstName)
		return "", fmt.Errorf("%s: %w", user.FirstName, ErrPandaNotFound)
	}

	return nick, nil
}

// TGUserIDFromNick return the Telegram user ID from a given nick
func TGUserIDFromNick(nick string) (int, error) {
	var TGUserID int

	query := fmt.Sprintf(`SELECT tg_user_id FROM user_nick WHERE nick = $1;`)
	err := pool.QueryRow(context.Background(), query, nick).Scan(TGUserID)
	if err != nil {
		log.Printf("Error finding user for %s. Fuck you.", nick)
		return NoneID, err
	}

	if nick == "" {
		log.Printf("%s doesn't have telegram.", nick)
		return NoneID, fmt.Errorf("%s: %w", nick, ErrPandaNotFound)
	}

	return TGUserID, nil
}
