package main

import (
	"sd_watcher/sd_app_config"
	"sd_watcher/telegram"
)

func main() {
	config := sd_app_config.CreateAppConfig("config.ini")
	checker := createChecker(config.SqlDataSourceName)
	if checker.isWorkingTimeNow() {
		checker.checkAppoitmentTimeout()
	}
	checker.checkPriceLoadTime(config.PriceParseTimeFile)
	if checker.msg != "" {
		telegram.SendMsg(checker.msg, config.TelegramApiToken, config.TelegramChatId)
	}
}
