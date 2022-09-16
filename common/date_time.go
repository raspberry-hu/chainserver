package common

import (
	"strconv"
	"time"
)

func IntToDate(timestamp int64) string {
	//timestamp, _ := strconv.ParseInt(str, 10, 64)
	timeobj := time.Unix(timestamp, 0)
	date := timeobj.Format("2006-01-02 15:04:05")
	return date
}

func StrToDate(str string) string {
	timestamp, _ := strconv.ParseInt(str, 10, 64)
	timeobj := time.Unix(timestamp, 0)
	date := timeobj.Format("2006-01-02 15:04:05")
	return date
}

func FixDate(value string) string {
	// value "2021-10-21T15:14:15Z"
	to, _ := time.Parse("2006-01-02T15:04:05Z", value)
	stamp := to.Format("2006-01-02 15:04:05")
	return stamp
}
