package config

import "time"

// 外部に公開するconst
const (
	ApproachURL                  = "https://transfer-cloud.navitime.biz/ohmitetudo/approachings"
	TimeTableURL                 = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	FrRits                       = "立命館大学〔近江鉄道・湖国バス〕"
	FrMinakusa                   = "南草津駅〔近江鉄道・湖国バス〕"
	TimeTableCacheUpdateDuration = 24 * 60 * 60 * time.Second
)

// 南草津タイムテーブルURL
var minamikusatsuTimeTableURL = map[string]string{
	"シャトル":          "https://transfer-cloud.navitime.biz/ohmitetudo/courses/timetables?busstop=00480156&course-sequence=0007900570-1",
	"かがやき通り":        "https://transfer-cloud.navitime.biz/ohmitetudo/courses/timetables?busstop=00480156&course-sequence=0007900505-1",
	"パナソニック西口・東口経由": "https://transfer-cloud.navitime.biz/ohmitetudo/courses/timetables?busstop=00480156&course-sequence=0007900503-1",
}

var ritsumeikanTimeTableURL = map[string]string{
	"invalid": "https://transfer-cloud.navitime.biz/ohmitetudo/courses/timetables?busstop=00480011&course-sequence=0007900502-1",
}

type From int

const (
	Unknown From = iota
	FromRits
	FromMinakusa
	FromNoji
	FromNandayama
	FromTamagawashogakkomae
	FromOnoyama
	FromPanaHigashi
	FromPanaMae
	FromPanaNishi
	FromKasayamaHigashi
	FromSasanoguchi
	FromKuresutoKusatsumae
	FromBkcGreenField
	FromNojiKigaguchi
	FromKusatsuKureaHole
	FromHigashiyakuraMinami
	FromHigashiyakuraShokuinnjutaku
	FromMukoyamaNewTown
	FromMaruo
	FromWakakusaKitaguchi
	FromRitsumeikanUnivMae
)

// String は From 型を文字列に変換します
func (f From) String() string {
	switch f {
	case Unknown:
		return "Unknown"
	case FromRits:
		return "立命館大学"
	case FromMinakusa:
		return "南草津駅"
	case FromNoji:
		return "野路"
	case FromNandayama:
		return "南田山"
	case FromTamagawashogakkomae:
		return "玉川小学校前"
	case FromOnoyama:
		return "小野山"
	case FromPanaHigashi:
		return "パナソニック東口"
	case FromPanaMae:
		return "パナソニック前"
	case FromPanaNishi:
		return "パナソニック西口"
	case FromKasayamaHigashi:
		return "笠山東"
	case FromSasanoguchi:
		return "笹の口"
	case FromKuresutoKusatsumae:
		return "クレスト草津前"
	case FromBkcGreenField:
		return "BKCグリーンフィールド"
	case FromNojiKigaguchi:
		return "野路北口"
	case FromKusatsuKureaHole:
		return "草津クレアホール"
	case FromHigashiyakuraMinami:
		return "東矢倉南"
	case FromHigashiyakuraShokuinnjutaku:
		return "東矢倉職員住宅"
	case FromMukoyamaNewTown:
		return "向山ニュータウン"
	case FromMaruo:
		return "丸尾"
	case FromWakakusaKitaguchi:
		return "若草北口"
	case FromRitsumeikanUnivMae:
		return "立命館大学正門前"
	default:
		return "Unknown"
	}
}

type To int

const (
	ToUnknown To = iota
	ToRits
	ToMinakusa
)

// String は To 型を文字列に変換します
func (t To) String() string {
	switch t {
	case ToUnknown:
		return "Unknown"
	case ToRits:
		return "立命館大学"
	case ToMinakusa:
		return "南草津駅"
	default:
		return "Unknown"
	}
}

// privateなconst
const (
	rits     = "立命館大学"
	minakusa = "南草津駅"
)

// frクエリに対応するdgmplクエリのスライス
var dgmplMap = map[string][]string{FrMinakusa: []string{"南草津駅〔近江鉄道・湖国バス〕:1", "南草津駅〔近江鉄道・湖国バス〕:3", "南草津駅〔近江鉄道・湖国バス〕:4"},
	FrRits: []string{"立命館大学〔近江鉄道・湖国バス〕:2"}}

// frクエリからスクレイピングURL用のfrに変換するマップ
var frList = map[string]string{
	rits:           "00480011",
	minakusa:       "00480156",
	"野路":           "00480081",
	"南田山":          "00480077",
	"玉川小学校前":       "00480047",
	"小野山":          "00480051",
	"パナソニック東口":     "00480078",
	"パナソニック前":      "00480079",
	"パナソニック西口":     "00480142",
	"笠山東":          "00481049",
	"笹の口":          "00480166",
	"クレスト草津前":      "00480167",
	"BKCグリーンフィールド": "00480168",
	"野路北口":         "00480082",
	"草津クレアホール":     "00480158",
	"東矢倉南":         "00480160",
	"東矢倉職員住宅":      "00480161",
	"向山ニュータウン":     "00480162",
	"丸尾":           "00480163",
	"若草北口":         "00480204",
	"立命館大学正門前":     "00480012",
}

// frクエリとtoクエリからdgmplクエリを取り出すマップ
var dgmplList = map[string]map[string][]string{
	rits:           {minakusa: {"立命館大学〔近江鉄道・湖国バス〕:2"}},
	minakusa:       {rits: {"南草津駅〔近江鉄道・湖国バス〕:1", "南草津駅〔近江鉄道・湖国バス〕:3", "南草津駅〔近江鉄道・湖国バス〕:4"}},
	"野路":           {rits: {"野路〔近江鉄道・湖国バス〕:1"}, minakusa: {"野路〔近江鉄道・湖国バス〕:2"}},
	"南田山":          {rits: {"南田山〔近江鉄道・湖国バス〕:1"}, minakusa: {"南田山〔近江鉄道・湖国バス〕:2"}},
	"玉川小学校前":       {rits: {"玉川小学校前〔近江鉄道・湖国バス〕:1"}, minakusa: {"玉川小学校前〔近江鉄道・湖国バス〕:2"}},
	"小野山":          {rits: {"小野山〔近江鉄道・湖国バス〕:1"}, minakusa: {"小野山〔近江鉄道・湖国バス〕:2"}},
	"パナソニック東口":     {rits: {"パナソニック東口〔近江鉄道・湖国バス〕:1", "パナソニック東口〔近江鉄道・湖国バス〕:3"}, minakusa: {"パナソニック東口〔近江鉄道・湖国バス〕:2"}},
	"パナソニック前":      {rits: {"パナソニック前〔近江鉄道・湖国バス〕:1", "パナソニック前〔近江鉄道・湖国バス〕:2"}},
	"パナソニック西口":     {rits: {"パナソニック西口〔近江鉄道・湖国バス〕:1", "パナソニック西口〔近江鉄道・湖国バス〕:2"}},
	"笠山東":          {rits: {"笠山東〔近江鉄道・湖国バス〕:2"}, minakusa: {"笠山東〔近江鉄道・湖国バス〕:1"}},
	"笹の口":          {rits: {"笹の口〔近江鉄道・湖国バス〕:1"}, minakusa: {"笹の口〔近江鉄道・湖国バス〕:2"}},
	"クレスト草津前":      {rits: {"クレスト草津前〔近江鉄道・湖国バス〕:1"}, minakusa: {"クレスト草津前〔近江鉄道・湖国バス〕:2"}},
	"BKCグリーンフィールド": {rits: {"ＢＫＣグリーンフィールド〔近江鉄道・湖国バス〕:2"}, minakusa: {"ＢＫＣグリーンフィールド〔近江鉄道・湖国バス〕:1"}},
	"野路北口":         {rits: {"野路北口〔近江鉄道・湖国バス〕:1"}, minakusa: {"野路北口〔近江鉄道・湖国バス〕:2"}},
	"草津クレアホール":     {rits: {"草津クレアホール〔近江鉄道・湖国バス〕:1"}, minakusa: {"草津クレアホール〔近江鉄道・湖国バス〕:2"}},
	"東矢倉南":         {rits: {"東矢倉南〔近江鉄道・湖国バス〕:1"}, minakusa: {"東矢倉南〔近江鉄道・湖国バス〕:2"}},
	"東矢倉職員住宅":      {rits: {"東矢倉職員住宅〔近江鉄道・湖国バス〕:1"}, minakusa: {"東矢倉職員住宅〔近江鉄道・湖国バス〕:2"}},
	"向山ニュータウン":     {rits: {"向山ニュータウン〔近江鉄道・湖国バス〕:1"}, minakusa: {"向山ニュータウン〔近江鉄道・湖国バス〕:2"}},
	"丸尾":           {rits: {"丸尾〔近江鉄道・湖国バス〕:1"}, minakusa: {"丸尾〔近江鉄道・湖国バス〕:2"}},
	"若草北口":         {rits: {"若草北口〔近江鉄道・湖国バス〕:1"}, minakusa: {"若草北口〔近江鉄道・湖国バス〕:2"}},
	"立命館大学正門前":     {rits: {"立命館大学正門前〔近江鉄道・湖国バス〕:2"}, minakusa: {"立命館大学正門前〔近江鉄道・湖国バス〕:1"}},
}
