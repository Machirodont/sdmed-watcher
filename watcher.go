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

func main() {
	config := createAppConfig("config.ini")
	checker := createChecker(config.sqlDataSourceName)
	if checker.isWorkingTimeNow() {
		checker.checkAppoitmentTimeout()
	}
	checker.checkPriceLoadTime(config.priceParseTimeFile)
	if checker.msg != "" {
		telegramMsg(checker.msg, config.telegramApiToken, config.telegramChatId)
	}
}
