package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	ritsToMinamikusatsu  = "立命館→南草津駅"
	minamikusatsuToRitsu = "南草津駅→立命館"
)

var urls = map[string]string{ritsToMinamikusatsu: "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl?mode=1&fr=立命館大学〔近江鉄道・湖国バス〕&frsk=B&tosk=&dgmpl=立命館大学〔近江鉄道・湖国バス〕:2:2&p=0,8,10",
	minamikusatsuToRitsu: "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl?mode=1&fr=南草津駅〔近江鉄道・湖国バス〕&frsk=B&tosk=&dgmpl=南草津駅〔近江鉄道・湖国バス〕:1:0&p=0,8,10"}

var tdList = []string{".column_day1_t2", ".column_day2_t2", ".column_day3_t2"}

type timeTable struct {
	WeekDays map[int][]oneTimeTable `json:"weekdays"`
	Saturday map[int][]oneTimeTable `json:"saturday"`
	Holiday  map[int][]oneTimeTable `json:"holiday"`
}

type oneTimeTable struct {
	Via string `json:"via"`
	Min string `json:"min"`
}

func newTimeTable() timeTable {
	return timeTable{
		WeekDays: map[int][]oneTimeTable{},
		Saturday: map[int][]oneTimeTable{},
		Holiday:  map[int][]oneTimeTable{},
	}
}

func scrapeTimeTable(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("mode")
	switch query {
	case "0":
		ritsToMinakusa(w, r)
	case "1":
		minakusaToRits(w, r)
	}
}

func ritsToMinakusa(w http.ResponseWriter, r *http.Request) {
	timeTable := scrapeRitsToMinakusa()
	data, _ := json.Marshal(timeTable)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func minakusaToRits(w http.ResponseWriter, r *http.Request) {
	timeTable := scrapeMinakusaToRits()
	data, _ := json.Marshal(timeTable)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func scrapeRitsToMinakusa() timeTable {
	doc, err := goquery.NewDocument(urls[ritsToMinamikusatsu])
	if err != nil {
		fmt.Println("スクレイピングにエラーが発生:", err)
		return newTimeTable()
	}
	timeTable := scrapeTimeInfo(doc)
	return timeTable
}

func scrapeMinakusaToRits() timeTable {
	doc, err := goquery.NewDocument(urls[minamikusatsuToRitsu])
	if err != nil {
		fmt.Println("スクレイピングにエラーが発生:", err)
		return newTimeTable()
	}
	timeTable := scrapeTimeInfo(doc)
	return timeTable
}

func scrapeTimeInfo(doc *goquery.Document) timeTable {
	timeTable := newTimeTable()
	doc.Find(".time").Each(func(_ int, s *goquery.Selection) {
		hour, err := strconv.Atoi(s.Text())
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		row := s.Parent()
		weekdays := row.Find(tdList[0])
		saturday := row.Find(tdList[1])
		holidays := row.Find(tdList[2])
		// 平日の情報をスクレイピング
		weekdays.Each(func(_ int, s *goquery.Selection) {
			s.Find("li").Each(func(_ int, s *goquery.Selection) {
				via := s.Find(".legend").Find("span").Text()
				min := strings.TrimSpace(s.Text())
				split := strings.Fields(min)
				if len(split) == 1 {
					min = split[0]
				} else {
					min = split[len(split)-1]
				}
				timeTable.WeekDays[hour] = append(timeTable.WeekDays[hour], oneTimeTable{Via: via, Min: min})
			})
		})
		// 土曜の情報をスクレイピング
		saturday.Each(func(_ int, s *goquery.Selection) {
			s.Find("li").Each(func(_ int, s *goquery.Selection) {
				via := s.Find(".legend").Find("span").Text()
				min := strings.TrimSpace(s.Text())
				split := strings.Fields(min)
				if len(split) == 1 {
					min = split[0]
				} else {
					min = split[len(split)-1]
				}
				timeTable.Saturday[hour] = append(timeTable.Saturday[hour], oneTimeTable{Via: via, Min: min})
			})
		})
		// 休日の情報をスクレイピング
		holidays.Each(func(_ int, s *goquery.Selection) {
			s.Find("li").Each(func(_ int, s *goquery.Selection) {
				via := s.Find(".legend").Find("span").Text()
				min := strings.TrimSpace(s.Text())
				split := strings.Fields(min)
				if len(split) == 1 {
					min = split[0]
				} else {
					min = split[len(split)-1]
				}
				timeTable.Holiday[hour] = append(timeTable.Holiday[hour], oneTimeTable{Via: via, Min: min})
			})
		})
	})
	return timeTable
}
