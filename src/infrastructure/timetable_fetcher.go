package infrastructure

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
)

type TimetableFetcher struct {
}

var tdList = []string{"td.schedule.column_day1_t2", "td.schedule.column_day2_t2", "td.schedule.column_day3_t2"}

var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
var userAgent = "Busdes! server"

func (fetcher TimetableFetcher) FetchTimetable(from config.From, to config.To) domain.TimeTable {
	timetable := domain.CreateNewTimeTable()
	scrapeUrl := config.CreateTimeTableUrl(from, to)
	fmt.Println("scrapeURL: ", scrapeUrl)

	for via, url := range scrapeUrl {
		// io.Reader経由でドキュメントにパースする
		doc, err := fetchTimeTableDocument(url)
		if err != nil {
			log.Printf("goquery.NewDocumentFromReader failed because of %v", err)
			return timetable
		}

		// 新しいHTML構造に対応したスクレイピング関数を呼び出す
		weekdays, saturdays, holidays := fetchTimeTableNew(doc, via)

		// 曜日ごとの時刻表に追加
		for hour, busTimes := range weekdays {
			timetable.Weekdays[hour] = append(timetable.Weekdays[hour], busTimes...)
		}

		for hour, busTimes := range saturdays {
			timetable.Saturdays[hour] = append(timetable.Saturdays[hour], busTimes...)
		}

		for hour, busTimes := range holidays {
			timetable.Holidays[hour] = append(timetable.Holidays[hour], busTimes...)
		}
	}

	fmt.Println("fetchTimeTableDocumentは成功しました")
	timetable.SortOneBusTime()
	return timetable
}

func fetchTimeTableDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	// この辺りのHeaderを設定しないと403が返された
	req.Header.Add("accept", accept)
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return &goquery.Document{}, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return &goquery.Document{}, err
	}
	return doc, nil
}

// 新しいHTML構造からのスクレイピング
func fetchTimeTableNew(doc *goquery.Document, via string) (map[int][]domain.OneBusTime, map[int][]domain.OneBusTime, map[int][]domain.OneBusTime) {
	weekdays := make(map[int][]domain.OneBusTime)
	saturdays := make(map[int][]domain.OneBusTime)
	holidays := make(map[int][]domain.OneBusTime)

	// 時間帯（時）の行を見つける
	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		hourElem := row.Find("th div.py-1.text-lg")
		if hourElem.Length() == 0 {
			return
		}

		hourStr := strings.TrimSpace(hourElem.Text())
		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			return
		}

		// 平日のセル
		weekdayCell := row.Find("td").First()
		if weekdayCell.Length() > 0 {
			weekdays[hour] = parseTimeCell(weekdayCell, via)
		}

		// 土日祝日のセル
		holidayCell := row.Find("td").Last()
		if holidayCell.Length() > 0 {
			holidayTimes := parseTimeCell(holidayCell, via)
			// 土曜日と祝日に同じデータを使用（分け方がわからないため）
			saturdays[hour] = holidayTimes
			holidays[hour] = holidayTimes
		}
	})

	return weekdays, saturdays, holidays
}

// セル内の時刻情報をパースする
func parseTimeCell(cell *goquery.Selection, via string) []domain.OneBusTime {
	var busTimes []domain.OneBusTime

	cell.Find("li div a").Each(func(i int, timeLink *goquery.Selection) {
		minStr := strings.TrimSpace(timeLink.Text())

		viaText := ""
		// 経由情報（P, か, 西など）
		if via != "invalid" {
			viaText = via
		} else {
			cell.Find("sup").Each(func(j int, sup *goquery.Selection) {
				if j == i {
					viaText = config.GetViaFullName(strings.TrimSpace(sup.Text()))
				}
			})
		}
		// 数字だけを抽出
		re := regexp.MustCompile(`\d+`)
		min := re.FindString(minStr)

		if min != "" {
			busTime := domain.OneBusTime{
				Via:     viaText,
				Min:     min,
				BusStop: getBusStopFromVia(viaText),
			}
			busTimes = append(busTimes, busTime)
		}
	})

	return busTimes
}

// 経由情報から乗り場情報を取得
func getBusStopFromVia(via string) string {
	switch via {
	case "P":
		return "4番乗り場"
	case "西":
		return "4番乗り場"
	case "か":
		return "3番乗り場"
	case "直":
		return "2番乗り場"
	default:
		return "1番乗り場"
	}
}
