package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	_ "github.com/go-sql-driver/mysql"
)

type SiteChecker struct {
	db                *sql.DB
	telegramApiToken  string
	telegramChatId    int64
	sqlDataSourceName string
	msg               string
}

func (checker *SiteChecker) check() {
	nanosecInHour := 60 * 60 * 1000000000
	alarmTime := time.Now().Add(time.Duration(-2 * nanosecInHour)).Format("2006-01-02 15:04:05")
	rows, err := checker.db.Query(`SELECT count(*) FROM sd_appointment WHERE status=0 AND created<"` + alarmTime + `"`)

	if err != nil {
		checker.msg = "Ошибка mysql"
	} else {
		for rows.Next() {
			var count int
			if err := rows.Scan(&count); err != nil {
				log.Fatal(err)
			}
			checker.msg = "Результатов запроса: " + fmt.Sprintf("%v", count)
		}
	}
}

func createChecker() SiteChecker {
	checker := SiteChecker{}
	checker.msg = ""

	configFileData, _ := ioutil.ReadFile("config.ini")
	configText := string(configFileData)
	re := regexp.MustCompile(`([a-zA-Z].*)=(.+)`)
	configLines := re.FindAllStringSubmatch(configText, -1)

	for _, line := range configLines {
		if line[1] == "sqlDataSourceName" {
			checker.sqlDataSourceName = strings.Trim(line[2], " \r\t")
		}
		if line[1] == "telegramApiToken" {
			checker.telegramApiToken = strings.Trim(line[2], " \r\t")
		}
		if line[1] == "telegramChatId" {
			checker.telegramChatId, _ = strconv.ParseInt(strings.Trim(line[2], " \r\t"), 10, 64)
		}
	}

	db, err := sql.Open("mysql", checker.sqlDataSourceName)
	if err != nil {
		fmt.Println("ERR")
	} else {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		err = db.Ping()

		if err != nil {
			fmt.Println("NOT CONNECTED")
		} else {
			fmt.Println("OK")
		}
	}
	checker.db = db
	return checker
}

func (checker *SiteChecker) telegramMsg(messageText string) {

	bot, err := tgbotapi.NewBotAPI(checker.telegramApiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	msg := tgbotapi.NewMessage(checker.telegramChatId, messageText)

	bot.Send(msg)

	updates, _ := bot.GetUpdatesChan(u)

	updates.Clear()
}

func main() {
	checker := createChecker()
	checker.msg = ""
	checker.check()
	checker.telegramMsg(checker.msg)
}
