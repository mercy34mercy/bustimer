package infrastructure

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
	"log"
	"net/http"
	"regexp"
)

type TimetableFetcher struct {
}

var tdList = []string{"td.schedule.column_day1_t2", "td.schedule.column_day2_t2", "td.schedule.column_day3_t2"}

var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
var userAgent = "Busdes! server"

func (fetcher TimetableFetcher) FetchTimetable(from config.From, to config.To) domain.TimeTable {
	timetable := domain.CreateNewTimeTable()
	scrapeUrl := config.CreateTimeTableUrl(from, to)
	fmt.Println("scrapeURL: ",scrapeUrl)

	// io.Reader経由でドキュメントにパースする
	doc, err := fetchTimeTableDocument(scrapeUrl)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader failed because of %v", err)
		return timetable
	}
	fmt.Println("fetchTimeTableDocumentは成功しました")

	reg := regexp.MustCompile("[0-9]+")
	for i, v := range tdList {
		doc.Find(v).Each(func(j int, s *goquery.Selection) {
			s.Find(".ttList li").Each(func(_ int, t *goquery.Selection) {
				oneBusTime := domain.OneBusTime {
					Via: config.GetViaFullName(t.Find(".legend span").Text()),
					Min: reg.FindString(t.Text()),
					BusStop: config.GetBusStop(from, to),
				}
				if i == 0 {
					timetable.Weekdays[j+5] = append(timetable.Weekdays[j+5], oneBusTime)
				} else if i == 1 {
					timetable.Saturdays[j+5] = append(timetable.Saturdays[j+5], oneBusTime)
				} else if i == 2 {
					timetable.Holidays[j+5] = append(timetable.Holidays[j+5], oneBusTime)
				}
			})
		})
	}
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