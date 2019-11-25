package timetable

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/PuerkitoBio/goquery"
)

const (
	url      = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	cacheExt = ".json"
	bucket   = "analog-subset-179214-busdes"
)

var tdList = []string{".column_day1_t2", ".column_day2_t2", ".column_day3_t2"}
var dgmplMap = map[string][]string{"立命館大学〔近江鉄道・湖国バス〕": []string{"立命館大学〔近江鉄道・湖国バス〕:2"},
	"南草津駅〔近江鉄道・湖国バス〕": []string{"南草津駅〔近江鉄道・湖国バス〕:1", "南草津駅〔近江鉄道・湖国バス〕:3", "南草津駅〔近江鉄道・湖国バス〕:4"}}

var currentFromRitsVersion = 1
var currentFromMinakusaVersion = 1
var currentFromRitsHash = ""
var currentFromMinakusaHash = ""

type timeTable struct {
	Version  int                    `json:"version"`
	WeekDays map[int][]oneTimeTable `json:"weekdays"`
	Saturday map[int][]oneTimeTable `json:"saturday"`
	Holiday  map[int][]oneTimeTable `json:"holiday"`
}

type oneTimeTable struct {
	Via     string `json:"via"`
	Min     string `json:"min"`
	BusStop string `json:"bus_stop"`
}

func newTimeTable() timeTable {
	return timeTable{
		Version:  0,
		WeekDays: map[int][]oneTimeTable{},
		Saturday: map[int][]oneTimeTable{},
		Holiday:  map[int][]oneTimeTable{},
	}
}

// ScrapeTimeTable クエリで指定された出発地点->目的地へ向かう時刻表データを返す
func ScrapeTimeTable(c echo.Context) error {
	fr := c.QueryParam("fr")
	// clientHash := c.QueryParam("hash")
	version, _ := strconv.Atoi(c.QueryParam("version"))
	// サーバ側でクラウドストレージによるキャッシュを用いて高速化する場合
	// fileName := fr + cacheExt
	// data, err := fetchFromCloudStorage(fileName)
	// if err == nil {
	// 	cacheHash, err := md5HashFromData(data)
	// 	if err == nil {
	// 		if clientHash == cacheHash || version == "1" {
	// 			return c.JSONBlob(http.StatusOK, data)
	// 		}
	// 	}
	// }

	// バージョンを確認して高速化する場合
	var currentVersion *int
	var currentHash *string
	if fr == "立命館大学〔近江鉄道・湖国バス〕" {
		currentVersion = &currentFromRitsVersion
		currentHash = &currentFromRitsHash
	} else if fr == "南草津駅〔近江鉄道・湖国バス〕" {
		currentVersion = &currentFromMinakusaVersion
		currentHash = &currentFromMinakusaHash
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "クエリfrを確認してください"})
	}
	if version == *currentVersion {
		return c.JSON(http.StatusNotModified, map[string]string{"reuslt": "Not modified"})
	}
	timeTable := newTimeTable()
	for _, v := range dgmplMap[fr] {
		fullURL := url + "?fr=" + fr + "&dgmpl=" + v
		err := scrapeFromURL(fullURL, &timeTable, string(v[len(v)-1]))
		if err != nil {
			return err
		}
		c.Echo().Logger.Debug("Successfully scrape from " + fullURL)
	}
	byteData, _ := json.Marshal(timeTable)
	newHash, _ := md5HashFromData(byteData)
	c.Echo().Logger.Debug("newHashData: " + newHash)
	c.Echo().Logger.Debug("currentHash: " + *currentHash)
	if newHash != *currentHash {
		*currentVersion++
		*currentHash = newHash
	}
	timeTable.Version = *currentVersion
	// err := saveCache(timeTable, fileName)
	// if err != nil {
	// return err
	// }
	// c.Echo().Logger.Debug("Successfully save file " + fileName)
	// c.Echo().Logger.Debug("Version: " + timeTable.Version)
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

func scrapeFromURL(fullURL string, timeTable *timeTable, busStop string) error {
	doc, err := goquery.NewDocument(fullURL)
	if err != nil {
		return err
	}
	scrapeTimeInfo(doc, timeTable, busStop)
	return nil
}

func scrapeTimeInfo(doc *goquery.Document, timeTable *timeTable, busStop string) {
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
				timeTable.WeekDays[hour] = append(timeTable.WeekDays[hour], oneTimeTable{Via: via, Min: min, BusStop: busStop})
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
				timeTable.Saturday[hour] = append(timeTable.Saturday[hour], oneTimeTable{Via: via, Min: min, BusStop: busStop})
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
				timeTable.Holiday[hour] = append(timeTable.Holiday[hour], oneTimeTable{Via: via, Min: min, BusStop: busStop})
			})
		})
	})
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
