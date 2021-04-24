package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func patron(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	const patronImagePath string = "/resources/images/patron.jpeg"
	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, patronImagePath)
	_, err := bot.Send(msg)
	if err != nil {
		msgErr := tgbotapi.NewMessage(update.Message.Chat.ID, "Error reading resource "+patronImagePath)
		bot.Send(msgErr)
	}
}
