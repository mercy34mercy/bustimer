package timetable

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo"

	"github.com/PuerkitoBio/goquery"
)

const (
	ritsToMinamikusatsu  = "立命館→南草津駅"
	minamikusatsuToRitsu = "南草津駅→立命館"
	url                  = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
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

// ScrapeTimeTable クエリで指定された出発地点->目的地へ向かう時刻表データを返す
func ScrapeTimeTable(c echo.Context) error {
	fr := c.QueryParam("fr")
	dgmpl := c.QueryParam("dgmpl")
}

func ritsToMinakusa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientHash := r.URL.Query().Get("hash")
	var hashString = ""
	file, err := os.Open("ritsToMinakusa.json")
	if err == nil {
		byteData, err := ioutil.ReadAll(file)
		if err == nil {
			hash := md5.New()
			file, err := os.Open("ritsToMinakusa.json")
			if err == nil {
				if _, err := io.Copy(hash, file); err == nil {
					hashInBytes := hash.Sum(nil)
					hashString = hex.EncodeToString(hashInBytes)
					fmt.Println(hashString)
				}
				if hashString == clientHash {
					w.Write(byteData)
					return
				}
			}
		}
	}
	timeTable := scrapeRitsToMinakusa()
	data, err := json.Marshal(timeTable)
	err = ioutil.WriteFile("ritsToMinakusa.json", data, 0644)
	if err != nil {
		fmt.Println("File save error:", err)
	}
	w.Write(data)
}

func minakusaToRits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	clientHash := r.URL.Query().Get("hash")
	var hashString = ""
	file, err := os.Open("minakusaToRits.json")
	if err == nil {
		byteData, err := ioutil.ReadAll(file)
		if err == nil {
			hash := md5.New()
			file, err := os.Open("minakusaToRits.json")
			if err == nil {
				if _, err := io.Copy(hash, file); err == nil {
					hashInBytes := hash.Sum(nil)
					hashString = hex.EncodeToString(hashInBytes)
					fmt.Println(hashString)
				}
				if hashString == clientHash {
					w.Write(byteData)
					return
				}
			}
		}
	}
	timeTable := scrapeMinakusaToRits()
	data, _ := json.Marshal(timeTable)
	err = ioutil.WriteFile("minakusaToRits.json", data, 0644)
	if err != nil {
		fmt.Println("File save error:", err)
	}
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
