package config

import (
	"github.com/yut-kt/goholiday"
	"log"
	"regexp"
	"time"
)

// frクエリからスクレイピングする際のfrクエリに変換する
// 該当するものがない場合は空文字列を返す
func GetFromKey(fr string) string {
	from, ok := frList[fr]
	if ok == false {
		log.Printf("fr=%v に対応するfromが設定されていません", fr)
		return ""
	}
	return from
}

// fr, toクエリからスクレイピングする際のdgmplクエリのスライスに変換する
// 該当するものがない場合は空のスライスを返す
func GetDgmplList(fr, to string) []string {
	dgmplMap, ok := dgmplList[fr]
	if ok == false {
		log.Printf("fr=%v に対応するdgmplMapが設定されていません", fr)
		return make([]string, 0)
	}
	dgmplList, ok := dgmplMap[to]
	if ok == false {
		log.Printf("to=%v に対応するdgmplListが設定されていません", to)
		return make([]string, 0)
	}
	return dgmplList
}

func CreateApproachInfoUrl(from, dgmpl string) string {
	return ApproachURL + "?fr=" + from + "&dgmpl=" + dgmpl
}

func CreateTimeTableUrl(from From, to To) string {
	fr := frList[from.toString()]
	dgmpl := dgmplList[from.toString()][to.toString()][0]
	return TimeTableURL + "?fr=" + fr + "&dgmpl=" + dgmpl
}

func ConvertFromFrQuery(fr string) From {
	switch fr {
	case rits:
		return FromRits
	case minakusa:
		return FromMinakusa
	case "野路":
		return FromNoji
	case "南田山":
		return FromNandayama
	case "玉川小学校前":
		return FromTamagawashogakkomae
	case "小野山":
		return FromOnoyama
	case "パナソニック東口":
		return FromPanaHigashi
	case "パナソニック前":
		return FromPanaMae
	case "パナソニック西口":
		return FromPanaNishi
	case "笠山東":
		return FromKasayamaHigashi
	case "笹の口":
		return FromSasanoguchi
	case "クレスト草津前":
		return FromKuresutoKusatsumae
	case "BKCグリーンフィールド":
		return FromBkcGreenField
	case "野路北口":
		return FromNojiKigaguchi
	case "草津クレアホール":
		return FromKusatsuKureaHole
	case "東矢倉南":
		return FromHigashiyakuraMinami
	case "東矢倉職員住宅":
		return FromHigashiyakuraShokuinnjutaku
	case "向山ニュータウン":
		return FromMukoyamaNewTown
	case "丸尾":
		return FromMaruo
	case "若草北口":
		return FromWakakusaKitaguchi
	case "立命館大学正門前":
		return FromRitsumeikanUnivMae
	default:
		log.Printf("fr: %v に該当するVIA_FROMは設定されていません", fr)
		return Unknown
	}
}

func ConvertFromToQuery(to string) To {
	switch to {
	case rits:
		return ToRits
	case minakusa:
		return ToMinakusa
	default:
		log.Printf("to: %v に該当するToは設定されていません", to)
		return ToUnknown
	}
}

func (from From) toString() string {
	switch from {
	case FromRits:
		return rits
	case FromMinakusa:
		return minakusa
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
	}
	return ""
}

func (to To) toString() string {
	switch to {
	case ToRits:
		return rits
	case ToMinakusa:
		return minakusa
	}
	return ""
}

func GetBusStop(from From, to To) string {
	dgmpl := dgmplList[from.toString()][to.toString()][0]
	reg := regexp.MustCompile(`[0-9]`)
	num := reg.FindString(dgmpl)
	return num
}

func GetVia(from From) string {
	switch from {
	case FromNoji, FromNandayama, FromTamagawashogakkomae, FromOnoyama, FromPanaHigashi:
		return GetViaFullName("P")
	case FromPanaMae, FromPanaNishi, FromKasayamaHigashi, FromSasanoguchi, FromKuresutoKusatsumae, FromBkcGreenField:
		return GetViaFullName("西")
	case FromNojiKigaguchi, FromKusatsuKureaHole, FromHigashiyakuraMinami, FromHigashiyakuraShokuinnjutaku, FromMukoyamaNewTown, FromMaruo, FromWakakusaKitaguchi, FromRitsumeikanUnivMae:
		return GetViaFullName("か")
	}
	return GetViaFullName("")
}

func IsHoliday() bool {
	datetime := time.Now()
	return goholiday.IsHoliday(datetime)
}

func GetViaFullName(via string) string {
	switch via {
	case "P":
		return "パナソニック東口経由"
	case "西":
		return "パナソニック西口経由"
	case "か":
		return "かがやき通り経由"
	case "立":
		return "立命館大学経由"
	case "P連":
		return "シャトルバス"
	default:
		return "シャトルバス"
	}
}

func TrimParentheses(str string) string {
	trimStr := str
	if trimStr[0] == '(' {
		trimStr = trimStr[1:]
	}
	if trimStr[len(trimStr) - 1] == ')' {
		trimStr = trimStr[:len(trimStr) - 1]
	}
	return trimStr
}

func GetRequiredTime(from From, to To, via string) int {
	switch to {
	case ToRits:
		switch from {
		case FromNoji, FromNandayama, FromTamagawashogakkomae:
			return 15
		case FromOnoyama, FromPanaHigashi:
			return 10
		case FromPanaMae, FromPanaNishi:
			return 15
		case FromKasayamaHigashi, FromSasanoguchi, FromKuresutoKusatsumae:
			return 10
		case FromBkcGreenField:
			return 5
		case FromNojiKigaguchi, FromKusatsuKureaHole:
			return 15
		case FromHigashiyakuraMinami, FromHigashiyakuraShokuinnjutaku, FromMukoyamaNewTown, FromMaruo:
			return 10
		case FromWakakusaKitaguchi, FromRitsumeikanUnivMae:
			return 5
		}
	case ToMinakusa:
		switch from {
		case FromNoji, FromNandayama, FromTamagawashogakkomae:
			return 5
		case FromOnoyama, FromPanaHigashi:
			return 10
		case FromPanaMae, FromPanaNishi:
			return 15
		case FromKasayamaHigashi, FromSasanoguchi, FromKuresutoKusatsumae:
			return 15
		case FromBkcGreenField:
			return 25
		case FromNojiKigaguchi, FromKusatsuKureaHole:
			return 5
		case FromHigashiyakuraMinami, FromHigashiyakuraShokuinnjutaku, FromMukoyamaNewTown, FromMaruo:
			return 10
		case FromWakakusaKitaguchi, FromRitsumeikanUnivMae:
			return 15
		}
	}
	if via == "パナソニック東口経由" {
		return 20
	} else if via == "パナソニック西口経由" {
		return 25
	} else if via == "かがやき通り経由" {
		return 20
	} else {
		return 15
	}
}