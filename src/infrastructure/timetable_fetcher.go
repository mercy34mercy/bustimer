package infrastructure

import (
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

func (fetcher TimetableFetcher) FetchTimetable(from string) domain.TimeTable {
	timetable := domain.CreateNewTimeTable()
	scrapeUrl := config.CreateTimeTableUrl(from)

	// ページの情報を取得する
	resp, err := http.Get(scrapeUrl)
	if err != nil {
		log.Printf("Http/GET failed to %v because of %v", scrapeUrl, err)
		return timetable
	}
	defer resp.Body.Close()

	// io.Reader経由でドキュメントにパースする
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader failed because of %v", err)
		return timetable
	}

	reg := regexp.MustCompile("[0-9]+")
	doc.Find(tdList[0]).Each(func(i int, s *goquery.Selection) {
		timetable.Weekdays[i + 5] = make([]domain.OneBusTime, 0)
		s.Find(".ttList li").Each(func(j int, t *goquery.Selection) {
			oneBusTime := domain.OneBusTime{
				Via: t.Find(".legend span").Text(),
				Min: reg.FindString(t.Text()),
				BusStop: getBusStop(from),
			}
			timetable.Weekdays[i + 5] = append(timetable.Weekdays[i + 5], oneBusTime)
		})
	})
	doc.Find(tdList[1]).Each(func(i int, s *goquery.Selection) {
		timetable.Saturdays[i + 5] = make([]domain.OneBusTime, 0)
		s.Find(".ttList li").Each(func(j int, t *goquery.Selection) {
			oneBusTime := domain.OneBusTime{
				Via: t.Find(".legend span").Text(),
				Min: reg.FindString(t.Text()),
				BusStop: getBusStop(from),
			}
			timetable.Saturdays[i + 5] = append(timetable.Saturdays[i + 5], oneBusTime)
		})
	})
	doc.Find(tdList[2]).Each(func(i int, s *goquery.Selection) {
		timetable.Holidays[i + 5] = make([]domain.OneBusTime, 0)
		s.Find(".ttList li").Each(func(j int, t *goquery.Selection) {
			oneBusTime := domain.OneBusTime{
				Via: t.Find(".legend span").Text(),
				Min: reg.FindString(t.Text()),
				BusStop: getBusStop(from),
			}
			timetable.Holidays[i + 5] = append(timetable.Holidays[i + 5], oneBusTime)
		})
	})
	return timetable
}

func getBusStop(from string) string {
	if from == config.FrRits {
		return "2"
	} else if from == config.FrMinakusa {
		return "1"
	} else {
		return ""
	}
}