package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
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
		checker.msg = "Ошибка соединения mysql [" + err.Error() + "]"
	} else {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		err = db.Ping()

		if err != nil {
			checker.msg = "Ошибка соединения mysql [" + err.Error() + "]"
		}
	}
	checker.db = db
	return checker
}

func (checker *SiteChecker) checkAppoitmentTimeout() {
	timeLimit := 2
	nanosecInHour := 60 * 60 * 1000000000
	alarmTime := time.Now().Add(time.Duration(-timeLimit * nanosecInHour)).Format("2006-01-02 15:04:05")
	rows, err := checker.db.Query(`SELECT count(*) FROM sd_appointment WHERE status=0 AND created<"` + alarmTime + `"`)

	if err != nil {
		checker.msg = "Ошибка mysql"
	} else {
		for rows.Next() {
			var count int
			if err := rows.Scan(&count); err != nil {
				log.Fatal(err)
			}
			if count > 0 {
				checker.msg += fmt.Sprintf("Больше %v-х часов висят необработанными заявки (%v шт.)\n", timeLimit, count)
			}
		}
	}
}

func (checker *SiteChecker) checkPriceLoadTime(filename string) {
	if _, err := os.Stat(filename); err == nil {
		fileData, _ := ioutil.ReadFile(filename)
		lastPriceLoad, _ := strconv.ParseInt(string(fileData), 10, 64)
		priceLoadInterval := time.Now().Unix() - lastPriceLoad
		secondsInWeek := 60 * 60 * 24 * 7
		if priceLoadInterval > int64(secondsInWeek) {
			formattedTime := time.Unix(lastPriceLoad, 0).Format("2006-01-02 15:04:05")
			checker.msg += "Цены не обновлялись больше недели (последний раз " + formattedTime + ") \n"
		}
	} else {
		checker.msg += "Нет данных о времени загрузки цен\n"
	}
}
