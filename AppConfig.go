package main

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type AppConfig struct {
	telegramApiToken   string
	telegramChatId     int64
	sqlDataSourceName  string
	priceParseTimeFile string
}

func createAppConfig(configFileName string) AppConfig {
	config := AppConfig{}
	configFileData, _ := ioutil.ReadFile(configFileName)
	configText := string(configFileData)
	re := regexp.MustCompile(`([a-zA-Z].*)=(.+)`)
	configLines := re.FindAllStringSubmatch(configText, -1)

	for _, line := range configLines {
		if line[1] == "sqlDataSourceName" {
			config.sqlDataSourceName = strings.Trim(line[2], " \r\t")
		}
		if line[1] == "telegramApiToken" {
			config.telegramApiToken = strings.Trim(line[2], " \r\t")
		}
		if line[1] == "telegramChatId" {
			config.telegramChatId, _ = strconv.ParseInt(strings.Trim(line[2], " \r\t"), 10, 64)
		}
		if line[1] == "priceParseTimeFile" {
			config.priceParseTimeFile = strings.Trim(line[2], " \r\t")
		}
	}
	return config
}
