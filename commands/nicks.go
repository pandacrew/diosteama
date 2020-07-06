package commands

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
	"github.com/pandacrew-net/diosteama/format"
)

func soy(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	if len(argv) != 2 {
		reply = fmt.Sprintf("Dime quien eres: !soy TuNick")
	} else {
		err := database.SetNick(update.Message.From, argv[1])
		if err != nil {
			if errors.Is(err, database.ErrPandaExists) {
				reply = fmt.Sprintf("Tu ya eres")
			} else {
				reply = fmt.Sprintf("Algo no fue bien: %s", err)
			}
		} else {
			reply = fmt.Sprintf("Vale, a partir de ahora eres <code>%s</code>", argv[1])
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func quienes(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	if len(argv) != 2 {
		reply = fmt.Sprintf("¿Por quien preguntas?")
	} else {
		username, err := database.TGUserFromNick(argv[1])
		if err != nil {
			if errors.Is(err, database.ErrPandaNotFound) {
				reply = fmt.Sprintf("No sé de quien me hablas.")
			}
			reply = fmt.Sprintf("Algo no fue bien: %s", err)
		} else {
			reply = fmt.Sprintf("@%s es el panda anteriormente conocido como %s", username, argv[1])
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func checkAdmin(bot *tgbotapi.BotAPI, ChatID int64, user *tgbotapi.User) bool {
	member, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID: ChatID,
		UserID: user.ID,
	})
	fmt.Printf("\n%v\n", member)
	return err == nil && (member.IsAdministrator() || member.IsCreator())
}

func es(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string

	if !checkAdmin(bot, update.Message.Chat.ID, update.Message.From) {
		log.Printf("%s esta intentado cambiar a alguien sin permiso", update.Message.From.UserName)
		return
	}

	if len(argv) != 2 {
		reply = fmt.Sprintf("Dime quien quieres que sea respondiendo al luser: !es TuNick")
	} else {
		if update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.From == nil {
			reply = fmt.Sprintf("No se de quien hablas")
		} else {
			user := update.Message.ReplyToMessage.From
			err := database.AdminSetNick(user, argv[1])
			if err != nil {
				if errors.Is(err, database.ErrPandaExists) {
					reply = fmt.Sprintf("Tu ya eres")
				} else {
					reply = fmt.Sprintf("Algo no fue bien: %s", err)
				}
			} else {
				reply = fmt.Sprintf("Vale, a partir de ahora <code>%s</code> es <code>%s</code>", format.PrettyUser(user), argv[1])
			}
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
