package domain

// 時刻表テーブルを現したデータ構造
type TimeTable struct {
	Weekdays 	map[int][]OneBusTime `json:"weekdays"`
	Saturdays 	map[int][]OneBusTime `json:"saturdays"`
	Holidays 	map[int][]OneBusTime `json:"holidays"`
}

// バス一台の運行情報
type OneBusTime struct {
	Via 	string `json:"via"`
	Min 	string `json:"min"`
	BusStop string `json:"bus_stop"`
}

func CreateNewTimeTable() TimeTable {
	// 初期化
	timetable := TimeTable{
		Weekdays: make(map[int][]OneBusTime),
		Saturdays: make(map[int][]OneBusTime),
		Holidays: make(map[int][]OneBusTime),
	}

	// 時刻表にあるデータを埋める
	for i := 5; i <= 24; i++ {
		timetable.Weekdays[i] = make([]OneBusTime, 0)
		timetable.Saturdays[i] = make([]OneBusTime, 0)
		timetable.Holidays[i] = make([]OneBusTime, 0)
	}
	return timetable
}