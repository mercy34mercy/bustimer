package infrastructure

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/domain"
	"log"
	"regexp"
)

type TimetableFetcher struct {

}

const (
	TimeTableUrl        = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	FrRits     = "立命館大学〔近江鉄道・湖国バス〕"
	FrMinakusa = "南草津駅〔近江鉄道・湖国バス〕"
	dgmplRits = "立命館大学〔近江鉄道・湖国バス〕:2"
	dgmplMinakusa = "南草津駅〔近江鉄道・湖国バス〕:1"
)

var tdList = []string{"td.schedule.column_day1_t2", "td.schedule.column_day2_t2", "td.schedule.column_day3_t2"}


func (fetcher TimetableFetcher) FetchTimetable(from string) domain.TimeTable {
	timetable := domain.TimeTable{
		Weekdays: make(map[int][]domain.OneBusTime),
		Saturdays: make(map[int][]domain.OneBusTime),
		Holidays: make(map[int][]domain.OneBusTime),
	}
	dgmpl := dgmplRits
	if from == FrMinakusa {
		dgmpl = dgmplMinakusa
	}
	scrapeUrl := TimeTableUrl + "?fr=" + from + "&dgmpl=" + dgmpl
	log.Printf("TimeTable url: %v", scrapeUrl)
	doc, err := goquery.NewDocument(scrapeUrl)
	if err != nil {
		log.Fatalf("timeTableのスクレイピングに失敗しています: %v", err)
		return timetable
	}
	reg := regexp.MustCompile("[0-9]+")
	// TODO: busStopだけ正しく取れてない
	doc.Find(tdList[0]).Each(func(i int, s *goquery.Selection) {
		timetable.Weekdays[i + 5] = make([]domain.OneBusTime, 0)
		s.Find(".ttList li").Each(func(j int, t *goquery.Selection) {
			oneBusTime := domain.OneBusTime{
				Via: t.Find(".legend span").Text(),
				Min: reg.FindString(t.Text()),
				BusStop: "0",
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
				BusStop: "0",
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
				BusStop: "0",
			}
			timetable.Holidays[i + 5] = append(timetable.Holidays[i + 5], oneBusTime)
		})
	})
	return timetable
}