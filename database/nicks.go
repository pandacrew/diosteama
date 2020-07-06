package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const UsersTable = "users"

// ErrPandaNotFound error returned when a user or nick is not in the DB
var ErrPandaNotFound = errors.New("No such Panda")

// ErrPandaExists error returned when a user or nick is not in the DB
var ErrPandaExists = errors.New("This panda already exists")

// TODO: Hay que guardar toda la informacion del usuario (Ver types.User en la api de TG)
//    Las funciones "*From*" deberian devolver el usuario. El caller ya se las apañara

// NickFromTGUser return the nick for a given Telegram User
func NickFromTGUser(user *tgbotapi.User) (string, error) {
	var nick string

	query := fmt.Sprintf(`SELECT nick FROM %s WHERE tg_userid = $1;`, UsersTable)
	err := pool.QueryRow(context.Background(), query, user.ID).Scan(nick)
	if err != nil {
		log.Printf("Error finding panda %s. Fuck you.", user.String())
		return "", err
	}

	if nick == "" {
		log.Printf("%s is not a panda.", user.String())
		return "", fmt.Errorf("%s: %w", user.String(), ErrPandaNotFound)
	}

	return nick, nil
}

// NickFromTGUserName return the nick for a given Telegram User
func NickFromTGUserName(username string) (string, error) {
	var nick string

	query := fmt.Sprintf(`SELECT nick FROM %s WHERE tg_username = $1;`, UsersTable)
	err := pool.QueryRow(context.Background(), query, username).Scan(nick)
	if err != nil {
		log.Printf("Error finding panda %s. Fuck you.", username)
		return "", err
	}

	if nick == "" {
		log.Printf("%s is not a panda.", username)
		return "", fmt.Errorf("%s: %w", username, ErrPandaNotFound)
	}

	return nick, nil
}

// TGUserFromNick return the Telegram user ID from a given nick
func TGUserFromNick(nick string) (string, error) {
	var TGUser string

	query := fmt.Sprintf(`SELECT tg_username FROM %s WHERE nick = $1;`, UsersTable)
	err := pool.QueryRow(context.Background(), query, nick).Scan(&TGUser)
	if err != nil {
		log.Printf("Error finding user for %s. Fuck you.", nick)
		return "", err
	}

	if nick == "" {
		log.Printf("%s doesn't have telegram.", nick)
		return "", fmt.Errorf("%s: %w", nick, ErrPandaNotFound)
	}

	return TGUser, nil
}

// SetNick associates the telegram user and given nick
func SetNick(user *tgbotapi.User, nick string) error {
	var count int
	query := fmt.Sprintf(`SELECT count(*) from %s WHERE nick = $1 or tg_id = $2`, UsersTable)
	err := pool.QueryRow(context.Background(), query, nick, user.ID).Scan(&count)

	if err != nil {
		return fmt.Errorf("Something happened on DB: %s", err)
	}

	if count > 0 {
		return fmt.Errorf("Nick or User: %w", ErrPandaExists)
	}

	insert := fmt.Sprintf("INSERT INTO %s (nick, tg_id, tg_username) VALUES ($1, $2)", UsersTable)
	_, err = pool.Exec(context.Background(), insert, nick, user.ID, user.UserName)
	if err != nil {
		return fmt.Errorf("Could not add user: %s", err)
	}

	return nil
}
