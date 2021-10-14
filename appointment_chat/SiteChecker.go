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
	checker.msg = ""

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
	return checker
}

func (checker *SiteChecker) CheckNewAppointments() error {
	lastReportedIdFile := "last_appointment_id.txt"
	fileData, err := ioutil.ReadFile(lastReportedIdFile)
	var lastReportedAssignId int
	if err != nil {
		lastReportedAssignId = 0
	}

	lastReportedAssignId, _ = strconv.Atoi(string(fileData))

	err = ioutil.WriteFile(lastReportedIdFile, []byte(strconv.Itoa(lastReportedAssignId)), 0777)
	if err != nil {
		fmt.Println("Недоступен для записи файл " + lastReportedIdFile + err.Error())
		os.Exit(1)
	}

	rows, err := checker.db.Query(`SELECT sd_appointment.id AS id, sd_clinics.city AS clinic  FROM sd_appointment LEFT JOIN sd_clinics ON sd_appointment.clinic_id=sd_clinics.id WHERE sd_appointment.status=0 AND sd_appointment.id>` + strconv.Itoa(lastReportedAssignId))
	if err != nil {
		return &myError{"ERR: Ошибка запроса mysql [" + err.Error() + "]"}
	}

	newAppointments := make(map[string]int)
	for rows.Next() {
		var id int
		var clinic string
		err := rows.Scan(&id, &clinic)
		if err != nil {
			return &myError{"ERR: Ошибка получения результата mysql [" + err.Error() + "]"}
		}
		if _, ok := newAppointments[clinic]; !ok {
			newAppointments[clinic] = 0
		}
		newAppointments[clinic]++
		if lastReportedAssignId < id {
			lastReportedAssignId = id
		}
	}

	if len(newAppointments) > 0 {
		checker.msg += "Новые заявки: \n"
		for clinicName, appointmentNumber := range newAppointments {
			checker.msg += clinicName + ": " + strconv.Itoa(appointmentNumber) + "\n"
		}
	}
	ioutil.WriteFile(lastReportedIdFile, []byte(strconv.Itoa(lastReportedAssignId)), 0777)
	return nil
}

func (checker *SiteChecker) isWorkingTimeNow() bool {
	return time.Now().Hour() > 9 && time.Now().Hour() < 20
}

type myError struct {
	msg string
}

func (err *myError) Error() string {
	return err.msg
}
