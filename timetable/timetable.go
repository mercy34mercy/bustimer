package timetable

const (
	url        = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	cacheExt   = ".json"
	bucket     = "analog-subset-179214-busdes"
	frRits     = "立命館大学〔近江鉄道・湖国バス〕"
	frMinakusa = "南草津駅〔近江鉄道・湖国バス〕"
)

// 時刻表データスクレイピングする時に縦列を示すクラス名
var tdList = []string{".column_day1_t2", ".column_day2_t2", ".column_day3_t2"}

// frクエリに対応するdgmplクエリのスライス
var dgmplMap = map[string][]string{frRits: []string{"立命館大学〔近江鉄道・湖国バス〕:2"},
	frMinakusa: []string{"南草津駅〔近江鉄道・湖国バス〕:1", "南草津駅〔近江鉄道・湖国バス〕:3", "南草津駅〔近江鉄道・湖国バス〕:4"}}

// バージョン管理に用いる変数
var currentFromRitsVersion = 1
var currentFromMinakusaVersion = 1
var currentFromRitsHash = ""
var currentFromMinakusaHash = ""

// timeTable 時刻表データを表す構造体
type timeTable struct {
	Version  int                    `json:"version"`
	WeekDays map[int][]oneTimeTable `json:"weekdays"`
	Saturday map[int][]oneTimeTable `json:"saturday"`
	Holidays map[int][]oneTimeTable `json:"Holidays"`
}

// oneTimeTable 時刻表におけるバス1台に関する情報を表す構造体
type oneTimeTable struct {
	Via     string `json:"via"`
	Min     string `json:"min"`
	BusStop string `json:"bus_stop"`
}
