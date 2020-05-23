package test

import (
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"log"
	"testing"
)

func TestScrapeTimetableFromRits(t *testing.T) {
	fetcher := infrastructure.TimetableFetcher{}
	timetable := fetcher.FetchTimetable(infrastructure.FrRits)
	if len(timetable.Weekdays) == 0 {
		t.Fatal("Weekdaysの中身が空です。")
	}
	if timetable.Weekdays[7][0].Via != "P" {
		t.Fatalf("Weekdays 6:03の経由駅がPではなく%vになっています。", timetable.Weekdays[6][0].Via)
	}
	if timetable.Saturdays[7][0].Via != "西" {
		t.Fatalf("Saturdayの7:30の経由駅が西ではなく%vになっています。", timetable.Saturdays[7][0].Via)
	}
	if timetable.Weekdays[7][0].Min != "3" {
		t.Fatalf("Weekdaysの7:03発のバスが%v分発になっています", timetable.Weekdays[7][0].Min)
	}
	if timetable.Weekdays[14][3].Via != "か" {
		t.Fatalf("Wekdaysの14:25の経由駅が「か」ではなく%vになっています。", timetable.Weekdays[14][3].Via)
	}
	log.Printf("Weekdays, 14:25, via: %v, min: %v, busStop: %v", timetable.Weekdays[14][3].Via, timetable.Weekdays[14][3].Min, timetable.Weekdays[14][3].BusStop)
}