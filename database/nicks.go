package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	pgx "github.com/jackc/pgx/v4"
)

/*
Se necesita la siguiente tabla.

CREATE TABLE public.{UsersTable} (
    id SERIAL PRIMARY KEY,
    nick varchar(80) NOT NULL UNIQUE,
    tg_username varchar(120) UNIQUE,
    tg_id integer NOT NULL UNIQUE
);

*/

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

	query := fmt.Sprintf(`SELECT nick FROM %s WHERE tg_id = $1;`, UsersTable)
	err := pool.QueryRow(context.Background(), query, user.ID).Scan(&nick)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("@%s: %w", user.String(), ErrPandaNotFound)
		}
		return "", fmt.Errorf("Error finding panda %s(%d): %w", user.String(), user.ID, err)
	}

	if nick == "" {
		return "", fmt.Errorf("%s: %w", user.String(), ErrPandaNotFound)
	}

	return nick, nil
}

// NickFromTGUserName return the nick for a given Telegram User
func NickFromTGUserName(username string) (string, error) {
	var nick string

	query := fmt.Sprintf(`SELECT nick FROM %s WHERE tg_username = $1;`, UsersTable)
	err := pool.QueryRow(context.Background(), query, username).Scan(&nick)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("@%s: %w", username, ErrPandaNotFound)
		}
		return "", fmt.Errorf("Error finding panda for %s: %w", username, err)
	}

	return nick, nil
}

// TGUserFromNick return the Telegram user ID from a given nick
func TGUserFromNick(nick string) (string, error) {
	var TGUser string

	query := fmt.Sprintf(`SELECT tg_username FROM %s WHERE nick = $1;`, UsersTable)
	err := pool.QueryRow(context.Background(), query, nick).Scan(&TGUser)
	if err != nil {
		return "", fmt.Errorf("Error finding user for %s: %w", nick, err)
	}

	if nick == "" {
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

	insert := fmt.Sprintf("INSERT INTO %s (nick, tg_id, tg_username) VALUES ($1, $2, $3)", UsersTable)
	_, err = pool.Exec(context.Background(), insert, nick, user.ID, user.UserName)
	if err != nil {
		return fmt.Errorf("Could not add user: %s", err)
	}

	log.Printf("Añadido %s como %s", user.String(), nick)

	return nil
}
