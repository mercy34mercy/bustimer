package timetable

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/shun-shun123/bus-timer/hash"
)

// ScrapeTimeTable クエリで指定された出発地点->目的地へ向かう時刻表データを返す
func ScrapeTimeTable(c echo.Context) error {
	fr := c.QueryParam("fr")
	version, err := strconv.Atoi(c.QueryParam("version"))
	if err != nil {
		c.Echo().Logger.Info("query [version] is invalid: ", version)
		version = -1
	}

	// バージョンを確認して高速化する場合
	var currentVersion *int
	var currentHash *string
	if fr == frRits {
		currentVersion = &currentFromRitsVersion
		currentHash = &currentFromRitsHash
	} else if fr == frMinakusa {
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
			c.Echo().Logger.Error("Scrape Error: ", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"result": err.Error()})
		}
		c.Echo().Logger.Debug("Successfully scrape from " + fullURL)
	}
	byteData, _ := json.Marshal(timeTable)
	newHash, _ := hash.MD5HashFromData(byteData)
	c.Echo().Logger.Debug("newHashData: " + newHash)
	c.Echo().Logger.Debug("currentHash: " + *currentHash)
	if newHash != *currentHash {
		*currentVersion++
		*currentHash = newHash
	}
	timeTable.Version = *currentVersion
	return c.JSON(http.StatusOK, timeTable)
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
		Holidays := row.Find(tdList[2])
		// 平日の情報をスクレイピング
		weekdays.Each(func(_ int, s *goquery.Selection) {
			pickupDaysData(s, timeTable.WeekDays, busStop, hour)
		})
		// 土曜の情報をスクレイピング
		saturday.Each(func(_ int, s *goquery.Selection) {
			pickupDaysData(s, timeTable.Saturday, busStop, hour)
		})
		// 休日の情報をスクレイピング
		Holidays.Each(func(_ int, s *goquery.Selection) {
			pickupDaysData(s, timeTable.Holidays, busStop, hour)
		})
	})
}

// pickupDaysData 時刻表データを列ごとに取得して追加する
func pickupDaysData(s *goquery.Selection, data map[int][]oneTimeTable, busStop string, hour int) {
	s.Find("li").Each(func(_ int, s *goquery.Selection) {
		via := s.Find(".legend").Find("span").Text()
		min := strings.TrimSpace(s.Text())
		split := strings.Fields(min)
		if len(split) == 1 {
			min = split[0]
		} else {
			min = split[len(split)-1]
		}
		data[hour] = append(data[hour], oneTimeTable{Via: via, Min: min, BusStop: busStop})
	})
}
