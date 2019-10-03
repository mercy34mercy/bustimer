package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var url = "https://ohmitetudo-bus.jorudan.biz/busstatedtl"

const (
	ritsumeikan   = "立命館大学〔近江鉄道・湖国バス〕"
	minamikusatsu = "立命館大学〔近江鉄道・湖国バス〕:2"
)

func getTimeTable(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("fr")
	if from == "" {
		from = ritsumeikan
	}
	to := r.URL.Query().Get("dgmpl")
	if to == "" {
		to = minamikusatsu
	}
	fullURL := url + "?fr=" + from + "&dgmpl=" + to
	fmt.Println(fullURL)
	doc, err := goquery.NewDocument(fullURL)
	if err != nil {
		fmt.Println("Scraping Error: ", err)
	}
	var busInfo []busTime
	leaveDocuments := doc.Find("li .time")
	reg, err := regexp.Compile("[0-9][0-9]:[0-9][0-9]")
	leaveDocuments.Each(func(_ int, s *goquery.Selection) {
		match := reg.FindString(s.Text())
		leave, arrive := calcLeaveAndArriveTime(match, 20)
		fmt.Printf("%v -> %v\n", leave, arrive)
		busInfo = append(busInfo, newBusTime(leave, arrive, "", false, 180, "2"))
	})

	data, err := json.Marshal(busInfo[:len(busInfo)/2])
	if err != nil {
		fmt.Println("JSON marshal error:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("charset", "utf-16")
	w.Write(data)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.HandleFunc("/bus/time", getTimeTable)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/bus/timeにGETリクエストを送ってください"))
	})
	http.HandleFunc("/bus/timetable", scrapeTimeTable)
	http.ListenAndServe(":"+port, nil)
}
