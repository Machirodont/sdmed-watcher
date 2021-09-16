package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func telegramMsg(
	messageText string,
	telegramApiToken string,
	telegramChatId int64,
) {

	bot, err := tgbotapi.NewBotAPI(telegramApiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	msg := tgbotapi.NewMessage(telegramChatId, messageText)

	bot.Send(msg)

	updates, _ := bot.GetUpdatesChan(u)

	updates.Clear()
}

func main() {
	config := createAppConfig("config.ini")
	checker := createChecker(config.sqlDataSourceName)
	checker.checkAppoitmentTimeout()
	checker.checkPriceLoadTime(config.priceParseTimeFile)
	telegramMsg(checker.msg, config.telegramApiToken, config.telegramChatId)
}
