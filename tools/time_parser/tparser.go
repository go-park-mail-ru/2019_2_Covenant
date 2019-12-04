package time_parser

import (
	"strings"
	"time"
)

const (
	DATE = "2006-01-02"
)

func GetDuration(val string) string {
	start := strings.Index(val, "T")
	end := strings.Index(val, "Z")
	return val[start+1 : end]
}

func StringToTime(val string) time.Time {
	date, _ := time.Parse(DATE, val)
	return date
}
