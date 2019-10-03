package main

import (
	"fmt"
	"strconv"
	"time"
)

// calcLeaveAndArriveTime hh:mmの書式から時間と分を分けて取り出す
func calcLeaveAndArriveTime(t string, addDuration int) (string, string) {
	h, _ := strconv.Atoi(t[:2])
	m, _ := strconv.Atoi(t[3:])
	leave := time.Date(2019, 1, 1, h, m, 0, 0, time.UTC)
	arrive := leave.Add(time.Duration(addDuration) * time.Minute)
	return fmt.Sprintf("%02d:%02d", leave.Hour(), leave.Minute()), fmt.Sprintf("%02d:%02d", arrive.Hour(), arrive.Minute())
}
