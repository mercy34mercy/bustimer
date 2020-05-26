package config

import "time"

// 外部に公開するconst
const (
	ApproachURL  = "https://ohmitetudo-bus.jorudan.biz/busstatedtl"
	TimeTableURL = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	FrRits       = "立命館大学〔近江鉄道・湖国バス〕"
	FrMinakusa   = "南草津駅〔近江鉄道・湖国バス〕"
	TimeTableCacheUpdateDuration = 24 * 60 * 60 * time.Second
)

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

// privateなconst
const (
	rits = "立命館大学"
	minakusa = "南草津駅"
)

// frクエリに対応するdgmplクエリのスライス
var dgmplMap = map[string][]string{FrMinakusa: []string{"南草津駅〔近江鉄道・湖国バス〕:1", "南草津駅〔近江鉄道・湖国バス〕:3", "南草津駅〔近江鉄道・湖国バス〕:4"},
	FrRits: []string{"立命館大学〔近江鉄道・湖国バス〕:2"}}

// frクエリからスクレイピングURL用のfrに変換するマップ
var frList = map[string]string{
	rits:            "立命館大学〔近江鉄道・湖国バス〕",
	minakusa:       	   "南草津駅〔近江鉄道・湖国バス〕",
	"野路":           	   "野路〔近江鉄道・湖国バス〕",
	"南田山":         	   "南田山〔近江鉄道・湖国バス〕",
	"玉川小学校前": 		   "玉川小学校前〔近江鉄道・湖国バス〕",
	"小野山":       		   "小野山〔近江鉄道・湖国バス〕",
	"パナソニック東口": 	   "パナソニック東口〔近江鉄道・湖国バス〕",
	"パナソニック前":   	   "パナソニック前〔近江鉄道・湖国バス〕",
	"パナソニック西口":   	   "パナソニック西口〔近江鉄道・湖国バス〕",
	"笠山東":       	       "笠山東〔近江鉄道・湖国バス〕",
	"笹の口":        	   "笹の口〔近江鉄道・湖国バス〕",
	"クレスト草津前":     	   "クレスト草津前〔近江鉄道・湖国バス〕",
	"BKCグリーンフィールド":  "ＢＫＣグリーンフィールド〔近江鉄道・湖国バス〕",
	"野路北口":             "野路北口〔近江鉄道・湖国バス〕",
	"草津クレアホール":      "草津クレアホール〔近江鉄道・湖国バス〕",
	"東矢倉南":             "東矢倉南〔近江鉄道・湖国バス〕",
	"東矢倉職員住宅":        "東矢倉職員住宅〔近江鉄道・湖国バス〕",
	"向山ニュータウン":      "向山ニュータウン〔近江鉄道・湖国バス〕",
	"丸尾":                "丸尾〔近江鉄道・湖国バス〕",
	"若草北口":             "若草北口〔近江鉄道・湖国バス〕",
	"立命館大学正門前":      "立命館大学正門前〔近江鉄道・湖国バス〕",
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