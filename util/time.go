package util

import "time"

const default_format = "2006-01-02 15:04:05"

var wdays = [7]string{"日", "月", "火", "水", "木", "金", "土"}

func JpCurrentTIme() string {
	t := time.Now()
	return t.Format(default_format) + " (" + wdays[t.Weekday()] + ")"
}
