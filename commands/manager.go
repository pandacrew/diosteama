package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func manager(update tgbotapi.Update, bot *tgbotapi.BotAPI, argv []string) {
	const managerImagePath string = "/resources/images/manager.jpeg"
	msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, managerImagePath)
	_, err := bot.Send(msg)
	if err != nil {
		msgErr := tgbotapi.NewMessage(update.Message.Chat.ID, "Error reading resource "+managerImagePath)
		bot.Send(msgErr)
	}
}
