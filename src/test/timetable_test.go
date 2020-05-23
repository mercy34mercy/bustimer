package test

import (
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"testing"
)

func TestScrapeTimetable(t *testing.T) {
	fetcher := infrastructure.TimetableFetcher{}
	timetable := fetcher.FetchTimetable(infrastructure.FrRits)
	if len(timetable.Weekdays) == 0 {
		t.Fatal("Weekdaysの中身が空です。")
	}
	if timetable.Weekdays[7][0].Via != "P" {
		t.Fatalf("Weekdays 6:03の経由駅がPではなく%vになっています。", timetable.Weekdays[6][0].Via)
	}
}