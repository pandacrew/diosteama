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
	if len(argv) != 1 {
		reply = "Dime quien eres: !soy TuNick"
	} else {
		nick := argv[0]
		err := database.SetNick(update.Message.From, nick)
		if err != nil {
			if errors.Is(err, database.ErrPandaExists) {
				reply = "Tu ya eres"
			} else {
				reply = fmt.Sprintf("Algo no fue bien: %s", err)
			}
		} else {
			reply = fmt.Sprintf("Vale, a partir de ahora eres <code>%s</code>", nick)
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}

func quienes(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string

	log.Printf("%v\n", argv)

	if len(argv) != 1 {
		reply = "¿Por quien preguntas?"
	} else {
		term := argv[0]
		username, err := database.TGUserFromNick(term)
		if err == nil {
			reply = fmt.Sprintf("@%s es el panda anteriormente conocido como %s", username, term)
		} else {
			nick, err := database.NickFromTGUserName(term)
			if err == nil {
				reply = fmt.Sprintf("@%s es el panda anteriormente conocido como %s", term, nick)
			} else {
				log.Printf("Algo no fue bien: %s", err)
				reply = "No sé de quien me hablas."
			}
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

	if len(argv) != 1 {
		reply = "Dime quien quieres que sea respondiendo al luser: !es TuNick"
	} else {
		nick := argv[0]

		if update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.From == nil {
			reply = "No se de quien hablas"
		} else {
			user := update.Message.ReplyToMessage.From
			err := database.AdminSetNick(user, nick)
			if err != nil {
				if errors.Is(err, database.ErrPandaExists) {
					reply = "Tu ya eres"
				} else {
					reply = fmt.Sprintf("Algo no fue bien: %s", err)
				}
			} else {
				reply = fmt.Sprintf("Vale, a partir de ahora <code>%s</code> es <code>%s</code>", format.PrettyUser(user), nick)
			}
		}
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ParseMode = "html"
	bot.Send(msg)
}
