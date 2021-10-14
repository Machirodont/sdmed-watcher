package sd_app_config

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type AppConfig struct {
	TelegramApiToken   string
	TelegramChatId     int64
	SqlDataSourceName  string
	PriceParseTimeFile string
}

func CreateAppConfig(configFileName string) AppConfig {
	config := AppConfig{}
	configFileData, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatalln("Config file isn't available: " + err.Error())
	}
	configText := string(configFileData)
	re := regexp.MustCompile(`([a-zA-Z].*)=(.+)`)
	configLines := re.FindAllStringSubmatch(configText, -1)

	for _, line := range configLines {
		if line[1] == "sqlDataSourceName" {
			config.SqlDataSourceName = strings.Trim(line[2], " \r\t")
		}
		if line[1] == "telegramApiToken" {
			config.TelegramApiToken = strings.Trim(line[2], " \r\t")
		}
		if line[1] == "telegramChatId" {
			config.TelegramChatId, _ = strconv.ParseInt(strings.Trim(line[2], " \r\t"), 10, 64)
		}
		if line[1] == "priceParseTimeFile" {
			config.PriceParseTimeFile = strings.Trim(line[2], " \r\t")
		}
	}
	return config
}
