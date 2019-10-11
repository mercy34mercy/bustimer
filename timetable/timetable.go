package timetable

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo"

	"github.com/PuerkitoBio/goquery"
)

const (
	url      = "https://ohmitetudo-bus.jorudan.biz/diagrampoledtl"
	cacheExt = ".json"
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
	filePath := fr + dgmpl + cacheExt
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

func connectCloudStorage(data []byte, fileName string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("Client create error:", err)
	}
	bkt := client.Bucket("analog-subset-179214.appspot.com")
	attrs, err := bkt.Attrs(ctx)
	if err != nil {
		fmt.Println("Error attrs: ", err)
	}
	fmt.Printf("bucket %s created at %s, is located in %s with storage class %s\n", attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)
	obj := bkt.Object(fileName)
	w := obj.NewWriter(ctx)
	if _, err := w.Write(data); err != nil {
		fmt.Println("Write error: ", err)
	}
	if err := w.Close(); err != nil {
		fmt.Println("writer close error: ", err)
	}
	r, err := obj.NewReader(ctx)
	if err != nil {
		fmt.Println("Reader create error: ", err)
	}
	defer r.Close()
	if _, err = io.Copy(os.Stdout, r); err != nil {
		fmt.Println("read error:", err)
	}
}

func saveCache(data interface{}, fileName string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	connectCloudStorage(jsonData, fileName)
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
	// file, err := os.Open(filePath)
	// if err != nil {
	// return "", err
	// }
	// defer file.Close()
	// hash := md5.New()
	// if _, err := io.Copy(hash, file); err != nil {
	// 	return "", err
	// }
	// hashByte := hash.Sum(nil)
	// hashString := hex.EncodeToString(hashByte)
	// fmt.Println(hashString)
	// return hashString, nil
	return "hogehogehugahuga", nil
}
