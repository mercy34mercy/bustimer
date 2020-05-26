package config

import "github.com/google/martian/log"

// frクエリからスクレイピングする際のfrクエリに変換する
// 該当するものがない場合は空文字列を返す
func GetFromKey(fr string) string {
	from, ok := frList[fr]
	if ok == false {
		log.Errorf("fr=%v に対応するfromが設定されていません", fr)
		return ""
	}
	return from
}

// fr, toクエリからスクレイピングする際のdgmplクエリのスライスに変換する
// 該当するものがない場合は空のスライスを返す
func GetDgmplList(fr, to string) []string {
	dgmplMap, ok := dgmplList[fr]
	if ok == false {
		log.Errorf("fr=%v に対応するdgmplMapが設定されていません", fr)
		return make([]string, 0)
	}
	dgmplList, ok := dgmplMap[to]
	if ok == false {
		log.Errorf("to=%v に対応するdgmplListが設定されていません", to)
		return make([]string, 0)
	}
	return dgmplList
}

func CreateApproachInfoUrl(from, dgmpl string) string {
	return ApproachURL + "?fr=" + from + "&dgmpl=" + dgmpl
}

func CreateTimeTableUrl(from string) string {
	dgmpl := dgmplMap[FrRits][0]
	if from == FrMinakusa {
		dgmpl = dgmplMap[FrMinakusa][0]
	}
	return TimeTableURL + "?fr=" + from + "&dgmpl=" + dgmpl
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
		log.Errorf("fr: %v に該当するVIA_FROMは設定されていません", fr)
		return Unknown
	}
}

func (viaFrom From) ToVia() string {
	switch viaFrom {
	case FromRitsumeikanUnivMae:
		return ""
	}
	return ""
}

