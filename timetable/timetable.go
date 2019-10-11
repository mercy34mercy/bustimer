package timetable

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"

	"github.com/PuerkitoBio/goquery"
)

const (
	url      = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	cacheExt = ".json"
	bucket   = "analog-subset-179214.appspot.com"
)

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
	clientHash := c.QueryParam("hash")
	fileName := fr + dgmpl + cacheExt
	data, err := fetchFromCloudStorage(fileName)
	if err == nil {
		cacheHash, err := md5HashFromData(data)
		if err == nil {
			if clientHash == cacheHash {
				return c.JSONBlob(http.StatusOK, data)
			}
		}
	}
	fullURL := url + "?fr=" + fr + "&dgmpl=" + dgmpl
	timeTable, err := scrapeFromURL(fullURL)
	if err != nil {
		return err
	}
	c.Echo().Logger.Debug("Successfully scrape from " + fullURL)
	err = saveCache(timeTable, fileName)
	if err != nil {
		return err
	}
	c.Echo().Logger.Debug("Successfully save file " + fileName)
	return c.JSON(http.StatusOK, timeTable)
}

func saveCache(data interface{}, fileName string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = saveToCloudStorage(jsonData, fileName)
	if err != nil {
		return err
	}
	return nil
}

func scrapeFromURL(fullURL string) (timeTable, error) {
	doc, err := goquery.NewDocument(fullURL)
	if err != nil {
		return timeTable{}, err
	}
	timeTable := scrapeTimeInfo(doc)
	return timeTable, nil
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

// md5HashFromFile fileNameで示されるfileからmd5でハッシュ値を算出する
func md5HashFromFile(fileName string) (string, error) {
	data, err := fetchFromCloudStorage(fileName)
	if err != nil {
		return "", err
	}
	hashString, err := md5HashFromData(data)
	if err != nil {
		return "", err
	}
	return hashString, nil
}

// md5HashFromData byteデータからハッシュ値を算出する
func md5HashFromData(data []byte) (string, error) {
	hash := md5.New()
	if _, err := hash.Write(data); err != nil {
		return "", err
	}
	hashByte := hash.Sum(nil)
	hashString := hex.EncodeToString(hashByte)
	fmt.Println(hashString)
	return hashString, nil
}
