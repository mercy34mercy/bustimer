package test

import (
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"testing"
)

type date int

const (
	weekdays date = iota
	saturday
	holidays
)

var fetcher = infrastructure.TimetableFetcher{}

/*
	発　: 立命館大学
	行先: 南草津駅
    乗り場: 2
	の時刻表を取得して、いくつかの具体的なデータと照らし合わせて正しくスクレイピングできてるかのテスト
 */
func TestScrapeTimetableFromRits(t *testing.T) {
	timetable := fetcher.FetchTimetable(config.FrRits)
	everydayCheck(t, &timetable)

	oneBusCheck(t, &timetable, 7, 0, "P", "3", "2", weekdays)
	oneBusCheck(t, &timetable, 7, 0, "西", "30", "2", saturday)
	oneBusCheck(t, &timetable, 7, 0, "西", "30", "2", holidays)

	oneBusCheck(t, &timetable, 21, 2, "か", "10", "2", weekdays)
	oneBusCheck(t, &timetable, 21, 2, "西", "35", "2", saturday)
	oneBusCheck(t, &timetable, 21, 2, "西", "35", "2", holidays)
}

func TestScrapeTimetableFromMinakusa(t *testing.T) {
	timetable := fetcher.FetchTimetable(config.FrMinakusa)
	everydayCheck(t, &timetable)

	oneBusCheck(t, &timetable, 6, 0, "P", "57", "1", weekdays)
	oneBusCheck(t, &timetable, 7, 0, "か", "0", "1", saturday)
	oneBusCheck(t, &timetable, 7, 0, "か", "0", "1", holidays)

	oneBusCheck(t, &timetable, 21, 0, "P", "1", "1", weekdays)
	oneBusCheck(t, &timetable, 21, 0, "西", "10", "1", saturday)
	oneBusCheck(t, &timetable, 21, 0, "西", "10", "1", holidays)
}

// 平日・土曜・休日のデータを取得できているか、件数でチェックする
func everydayCheck(t *testing.T, timetable *domain.TimeTable) {
	if len(timetable.Weekdays) == 0 {
		t.Fatal("Weekdaysの中身が空です。")
	}
	if len(timetable.Saturdays) == 0 {
		t.Fatal("Saturdayの中身が空です。")
	}
	if len(timetable.Holidays) == 0 {
		t.Fatal("Holidaysの中身が空です。")
	}
}

// 特定の一つのバスに対する情報が全て正しいかチェックする
func oneBusCheck(t *testing.T, timetable *domain.TimeTable, hour int, index int, via, min, busStop string, kind date) {
	var target domain.OneBusTime
	switch kind {
	case weekdays:
		target = timetable.Weekdays[hour][index]
	case saturday:
		target = timetable.Saturdays[hour][index]
	case holidays:
		target = timetable.Holidays[hour][index]
	}
	if target.Via != via {
		t.Fatalf("%vの%v時発%v番目のバスのViaが%vになっています", kind, hour, index, target.Via)
	}
	if target.Min != min {
		t.Fatalf("%vの%v時発%v番目のバスのMinが%vになっています", kind, hour, index, target.Min)
	}
	if target.BusStop != busStop {
		t.Fatalf("%vの%v時発%v番目のバスのBusStopが%vになっています", kind, hour, index, target.BusStop)
	}
}