package infrastructure

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
	"github.com/shun-shun123/bus-timer/src/slack"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

func TimeTableRequest(c Context) error {
	from, _ := c.GetFromToQuery()
	timeTable, ok := TimeTableCache[from]
	if ok == false {
		if from == config.Unknown {
			return c.Response("TimeTableRequest", http.StatusBadRequest, timeTable)
		}
	}
	if len(timeTable.Weekdays[10]) <= 0 {
		slack.PostMessage(fmt.Sprint("従来の時刻表データのスクレイピング結果が不正なので、旧サイトからとります"))
		if from == config.FromRits {
			timeTable = TimeTableFetchFromRits()
		} else {
			timeTable = TimeTableFetchFromMinakusa()
			return c.Response("TimeTableRequest", http.StatusBadRequest, timeTable)
		}
		TimeTableCache[from] = timeTable
	}
	return c.Response("TimeTableRequest", http.StatusOK, timeTable)
}

func TimeTableFetchFromMinakusa() domain.TimeTable {
	timeTable := domain.CreateNewTimeTable()
	doc, _ := fetchDocument("http://time.khobho.co.jp/ohmi_bus/tim_dsp.asp?projCd=1&eigCd=7&teicd=1250&KaiKbn=NOW&pole=1")
	weekdays := fetchTimeTable(doc)
	doc, _ = fetchDocument("http://time.khobho.co.jp/ohmi_bus/tim_dsp.asp?projCd=2&eigCd=7&teicd=1250&KaiKbn=NOW&pole=1")
	saturdays := fetchTimeTable(doc)
	doc, _ = fetchDocument("http://time.khobho.co.jp/ohmi_bus/tim_dsp.asp?projCd=3&eigCd=7&teicd=1250&KaiKbn=NOW&pole=1")
	holidays := fetchTimeTable(doc)
	timeTable.Weekdays = weekdays
	timeTable.Saturdays = saturdays
	timeTable.Holidays = holidays
	return timeTable
}

func TimeTableFetchFromRits() domain.TimeTable {
	timeTable := domain.CreateNewTimeTable()
	doc, _ := fetchDocument("http://time.khobho.co.jp/ohmi_bus/tim_dsp.asp?projCd=1&eigCd=7&teicd=1050&KaiKbn=NOW&pole=2")
	weekdays := fetchTimeTableFromRits(doc)
	doc, _ = fetchDocument("http://time.khobho.co.jp/ohmi_bus/tim_dsp.asp?projCd=2&eigCd=7&teicd=1050&KaiKbn=NOW&pole=2")
	saturdays := fetchTimeTableFromRits(doc)
	doc, _ = fetchDocument("http://time.khobho.co.jp/ohmi_bus/tim_dsp.asp?projCd=3&eigCd=7&teicd=1050&KaiKbn=NOW&pole=2")
	holidays := fetchTimeTableFromRits(doc)
	timeTable.Weekdays = weekdays
	timeTable.Saturdays = saturdays
	timeTable.Holidays = holidays
	return timeTable
}

func fetchDocument(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return &goquery.Document{}, err
	}
	defer resp.Body.Close()

	// 読み取り
	buf, _ := ioutil.ReadAll(resp.Body)

	// 文字コード判定
	det := chardet.NewTextDetector()
	detResult, _ := det.DetectBest(buf)

	// 文字コード変換
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detResult.Charset, bReader)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return &goquery.Document{}, err
	}
	return doc, nil
}

func fetchTimeTable(doc *goquery.Document) map[int][]domain.OneBusTime {
	timeTable := map[int][]domain.OneBusTime{}
	regMin := regexp.MustCompile("[0-9]+")
	regVia := regexp.MustCompile("[^0-9]*")
	hours := doc.Find(".StyleTextCenter")
	hours.Each(func(index int, s *goquery.Selection) {
		hour, _ := strconv.Atoi(s.Text())
		min := ""
		via := ""
		mins := s.Parent().Find("nobr")
		mins.Each(func(i int, c *goquery.Selection) {
			min = regMin.FindString(c.Text())
			via = regVia.FindString(c.Text())
			oneBusTime := domain.OneBusTime{
				Via: config.GetViaFullName(via),
				Min: min,
				BusStop: getBusStopFromMinsakusa(via),
			}
			timeTable[hour] = append(timeTable[hour], oneBusTime)
		})
	})
	return timeTable
}

func fetchTimeTableFromRits(doc *goquery.Document) map[int][]domain.OneBusTime {
	timeTable := map[int][]domain.OneBusTime{}
	regMin := regexp.MustCompile("[0-9]+")
	regVia := regexp.MustCompile("[^0-9]*")
	hours := doc.Find(".StyleTextCenter")
	hours.Each(func(index int, s *goquery.Selection) {
		hour, _ := strconv.Atoi(s.Text())
		min := ""
		via := ""
		mins := s.Parent().Find("nobr")
		mins.Each(func(i int, c *goquery.Selection) {
			min = regMin.FindString(c.Text())
			via = regVia.FindString(c.Text())
			oneBusTime := domain.OneBusTime{
				Via: config.GetViaFullName(via),
				Min: min,
				BusStop: "1番乗り場",
			}
			timeTable[hour] = append(timeTable[hour], oneBusTime)
		})
	})
	return timeTable
}

func getBusStopFromMinsakusa(via string) string {
	switch via {
	case "P":
		return "4番乗り場"
	case "西":
		return "4番乗り場"
	case "か":
		return "3番乗り場"
	case "立":
		return "5番乗り場"
	default:
		return "1番乗り場"
	}
}