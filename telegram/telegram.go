package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendMsg(
	messageText string,
	telegramApiToken string,
	telegramChatId int64,
) {

	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		log.Fatalln("Telegram bot error: " + err.Error())
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	msg := tgbotapi.NewMessage(telegramChatId, messageText)

	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		log.Fatalln("Telegram msg sending error: " + sendErr.Error())
	}

	updates, _ := bot.GetUpdatesChan(u)

	updates.Clear()
}
