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
	url       = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	cachePath = "timetableCache"
	cacheExt  = ".json"
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
	filePath := cachePath + "/" + fr + dgmpl + cacheExt
	cacheHash, err := md5HashFromFile(filePath)
	if err == nil {
		if clientHash == cacheHash {
			file, err := os.Open(filePath)
			if err == nil {
				bytedata, err := ioutil.ReadAll(file)
				if err == nil {
					return c.JSONBlob(http.StatusOK, bytedata)
				}
			}
		}
	}
	fullURL := url + "?fr=" + fr + "&dgmpl=" + dgmpl
	timeTable, err := scrapeFromURL(fullURL)
	if err != nil {
		return err
	}
	c.Echo().Logger.Debug("Successfully scrape from " + fullURL)
	err = saveCache(timeTable, filePath)
	if err != nil {
		return err
	}
	c.Echo().Logger.Debug("Successfully save file " + filePath)
	return c.JSON(http.StatusOK, timeTable)
}

func saveCache(data interface{}, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
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

func md5HashFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashByte := hash.Sum(nil)
	hashString := hex.EncodeToString(hashByte)
	fmt.Println(hashString)
	return hashString, nil
}
