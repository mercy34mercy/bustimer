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

type To int

const (
	ToUnknown To = iota
	ToRits
	ToMinakusa
)

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
