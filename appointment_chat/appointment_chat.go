package main

import (
	"os"
	"sd_watcher/sd_app_config"
	"sd_watcher/telegram"
)

func main() {
	config := sd_app_config.CreateAppConfig("config.ini")
	checker := createChecker(config.SqlDataSourceName)

	if checker.db == nil {
		telegram.SendMsg(checker.msg, config.TelegramApiToken, config.TelegramChatId)
		os.Exit(1)
	}
	if checker.isWorkingTimeNow() {
		err := checker.CheckNewAppointments()
		if err == nil {
			if checker.msg != "" {
				telegram.SendMsg(checker.msg, config.TelegramApiToken, config.TelegramChatId)
			}
		} else {
			telegram.SendMsg(err.Error(), config.TelegramApiToken, config.TelegramChatId)
			os.Exit(1)
		}
	}

	if err := checker.db.Close(); err != nil {
		telegram.SendMsg(err.Error(), config.TelegramApiToken, config.TelegramChatId)
		os.Exit(1)
	}
}
