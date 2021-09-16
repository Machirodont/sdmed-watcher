package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SiteChecker struct {
	db  *sql.DB
	msg string
}

func createChecker(sqlDataSourceName string) SiteChecker {
	checker := SiteChecker{}
	checker.msg = "Тестовый режим\n"

	db, err := sql.Open("mysql", sqlDataSourceName)
	if err != nil {
		checker.msg += "ERR: Ошибка соединения mysql [" + err.Error() + "]\n"
	} else {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		err = db.Ping()

		if err != nil {
			checker.msg += "ERR: Ошибка соединения mysql [" + err.Error() + "]\n"
		} else {
			checker.db = db
		}
	}
	fmt.Println(checker.db)
	return checker
}

func (checker *SiteChecker) checkAppoitmentTimeout() {
	if checker.db == nil {
		return
	}
	timeLimit := 2
	nanosecInHour := 60 * 60 * 1000000000
	alarmTime := time.Now().Add(time.Duration(-timeLimit * nanosecInHour)).Format("2006-01-02 15:04:05")
	rows, err := checker.db.Query(`SELECT count(*) FROM sd_appointment WHERE status=0 AND created<"` + alarmTime + `"`)

	if err != nil {
		checker.msg += "ERR: Ошибка запроса mysql [" + err.Error() + "]\n"
	} else {
		for rows.Next() {
			var count int
			if err := rows.Scan(&count); err != nil {
				checker.msg += "ERR: Ошибка получения результата mysql [" + err.Error() + "]\n"
			}
			if count > 0 {
				checker.msg += fmt.Sprintf("Больше %v-х часов висят необработанными заявки (%v шт.)\n", timeLimit, count)
			}
		}
	}
}

func (checker *SiteChecker) checkPriceLoadTime(filename string) {
	if _, err := os.Stat(filename); err != nil {
		checker.msg += "ERR: Недоступен файл со временем последнего обновления цен\n"
		return
	}
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		checker.msg += "ERR: Ошибка чтения файла со временем последнего обновления цен\n"
		return
	}

	lastPriceLoad, err := strconv.ParseInt(string(fileData), 10, 64)
	if err != nil {
		checker.msg += "ERR: В файле со временем последнего обновления цен какая-то дичь\n"
		return
	}

	priceLoadInterval := time.Now().Unix() - lastPriceLoad
	secondsInWeek := 60 * 60 * 24 * 7
	if priceLoadInterval > int64(secondsInWeek) {
		formattedTime := time.Unix(lastPriceLoad, 0).Format("2006-01-02 15:04:05")
		checker.msg += "Цены не обновлялись больше недели (последний раз " + formattedTime + ") \n"
	}
}
