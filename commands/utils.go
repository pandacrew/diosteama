package commands

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func checkAdmin(bot *tgbotapi.BotAPI, ChatID int64, user *tgbotapi.User) bool {
	member, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
		ChatID: ChatID,
		UserID: user.ID,
	})
	fmt.Printf("\n%v\n", member)
	return err == nil && (member.IsAdministrator() || member.IsCreator())
}
