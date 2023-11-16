package util

import "time"

const TimeLayout = "2006-01-02T15:04:05"

var tz, _ = time.LoadLocation("Europe/Istanbul")

func Now() time.Time {
	return time.Now().In(tz)
}

func GetFormattedNow() string {
	return Now().Format(TimeLayout)
}
