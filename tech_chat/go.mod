module sd_watcher/watcher

go 1.17

replace sd_watcher/sd_app_config => ../sd_app_config

replace sd_watcher/telegram => ../telegram

require (
	github.com/go-sql-driver/mysql v1.6.0
	sd_watcher/sd_app_config v0.0.0-00010101000000-000000000000
	sd_watcher/telegram v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)
