package commands

import (
	"errors"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pandacrew-net/diosteama/database"
	"github.com/pandacrew-net/diosteama/format"
)

var SacredNames = map[string]string{
	"El Fary":       "Solo existe un auténtico Fary y no eres tu.",
	"una taza":      "Yo soy la tetera.",
	"Steve Ballmer": "DEVELOPERS DEVELOPERS DEVELOPERS!!!",
}

func soy(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	if len(argv) != 1 {
		reply = fmt.Sprintf("Dime quien eres: !soy TuNick")
	} else {
		nick := argv[0]
		err := database.SetNick(update.Message.From, nick)
		if err != nil {
			if errors.Is(err, database.ErrPandaExists) {
				reply = fmt.Sprintf("Tu ya eres")
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

func quienesNick(nick string) (string, error) {
	var reply = ""
	_, username, err := database.TGUserFromNick(nick)
	if err == nil {
		if username == "" {
			// Exists, but doesn't have a username
			reply = fmt.Sprintf("El panda anteriormente conocido como %s ahora es <pre>anonymous</pre>", nick)
		} else {
			reply = fmt.Sprintf("@%s es el panda anteriormente conocido como %s", username, nick)
		}
	}

	return reply, err
}

func quienesUsername(username string) (string, error) {
	var reply = ""
	if username[0] == '@' {
		username = username[1:]
	}

	nick, err := database.NickFromTGUserName(username)
	if err == nil {
		reply = fmt.Sprintf("@%s es el panda anteriormente conocido como %s", username, nick)
	}

	return reply, err
}

func quienesTGUser(user *tgbotapi.User) (string, error) {
	var reply = ""
	nick, err := database.NickFromTGUser(user)
	if err == nil {
		reply = fmt.Sprintf("@%s es el panda anteriormente conocido como %s", user.UserName, nick)
	}
	return reply, err
}

func quienes(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	var reply string
	var err error
	var exists bool

	if update.Message.ReplyToMessage != nil {
		reply, _ = quienesTGUser(update.Message.ReplyToMessage.From)
	} else if len(argv) < 1 {
		reply = fmt.Sprintf("¿Por quien preguntas?")
	} else {
		term := strings.Join(argv, " ")
		if reply, exists = SacredNames[term]; !exists {
			reply, err = quienesNick(term)
			if err != nil {
				reply, err = quienesUsername(term)
			}
		}
	}

	if reply == "" {
		reply = fmt.Sprintf("No sé de quien me hablas.")
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
		reply = fmt.Sprintf("Dime quien quieres que sea respondiendo al luser: !es TuNick")
	} else {
		nick := argv[0]

		if update.Message.ReplyToMessage == nil || update.Message.ReplyToMessage.From == nil {
			reply = fmt.Sprintf("No se de quien hablas")
		} else {
			user := update.Message.ReplyToMessage.From
			err := database.AdminSetNick(user, nick)
			if err != nil {
				if errors.Is(err, database.ErrPandaExists) {
					reply = fmt.Sprintf("Tu ya eres")
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
