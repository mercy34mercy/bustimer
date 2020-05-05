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