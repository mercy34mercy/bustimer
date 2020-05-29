package infrastructure

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type ApproachInfoFetcher struct {
	from config.From
	to config.To
}

type CustomDocument struct {
	*goquery.Document
}

var dataCount = 3

// 接近情報のサイトから取れる下記の情報をまとめてスクレイピングする
func (doc *CustomDocument) fetchApproachInfo() ([]string, []string, []string, []string, []string) {
	moreMin := make([]string, 0)
	realArrivalTime := make([]string, 0)
	Direction := make([]string, 0)
	ScheduledTime := make([]string, 0)
	Delay := make([]string, 0)
	doc.Find(".tableDetail").Each(func(i int, s *goquery.Selection) {
		if i >= dataCount {
			return
		}
		moreMin = append(moreMin, s.Find(".more_min").Text())
		r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
		arrivalTime := r.FindString(s.Find(".time").Text())
		realArrivalTime = append(realArrivalTime, arrivalTime)
		s.Find(".bsul li").Each(func(j int, li *goquery.Selection) {
			if j == 4 {
				trimed := strings.TrimSpace(strings.Split(li.Text(), "\n")[2])
				Direction = append(Direction, trimed)
			}
		})
		s.Find(".moreArea .bsul li").Each(func(j int, li *goquery.Selection) {
			r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
			scheduledTime := r.FindStringSubmatch(li.Text())
			if scheduledTime != nil {
				ScheduledTime = append(ScheduledTime, scheduledTime[0])
			}
			li.Find(".bsmidashi").Each(func(k int, span *goquery.Selection) {
				if k == 1 {
					Delay = append(Delay, config.TrimParentheses(span.Text()))
				}
			})
		})
	})
	return moreMin, realArrivalTime, Direction, ScheduledTime, Delay
}

func findMinLen(dataset ...[]string) int {
	min := 10000
	for _, v := range dataset {
		if len(v) < min {
			min = len(v)
		}
	}
	return min
}

func (fetcher ApproachInfoFetcher) FetchApproachInfos(approachInfoUrl string) domain.ApproachInfos {
	// 返り値で返す変数を初期化
	approachInfos := domain.CreateApproachInfos()

	// ページの情報を取得する
	resp, err := http.Get(approachInfoUrl)
	if err != nil {
		log.Printf("Http/GET failed to %v because of %v", approachInfoUrl, err)
		return approachInfos
	}
	defer resp.Body.Close()

	// io.Reader経由でドキュメントにパースする
	approachDoc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader failed because of %v", err)
		return approachInfos
	}

	// CustomDocument型に変換する
	customDoc := CustomDocument{approachDoc}
	moreMin, realArrivalTime, directions, scheduledTime, delay := customDoc.fetchApproachInfo()

	// どれかが空の場合もあるので、最小の数を探す
	iterateCount := findMinLen(moreMin, realArrivalTime, directions, scheduledTime, delay)
	via := ""
	for i := 0; i < iterateCount; i++ {
		//TODO: 経由情報のスクレイピング

		// hh:mmの表記でくる
		hour, _ := strconv.Atoi(scheduledTime[i][:2])
		min, _ := strconv.Atoi(scheduledTime[i][3:])
		tt, ok := TimeTable[fetcher.from]
		if ok {
			timeTableData := tt.Weekdays
			if config.IsHoliday() {
				timeTableData = tt.Saturdays
			}
			for _, v := range timeTableData[hour] {
				if convMin, err := strconv.Atoi(v.Min); err == nil {
					if convMin == min {
						via = v.Via
					}
				}
			}
		} else if fetcher.from != config.Unknown {
			via = config.GetVia(fetcher.from)
		}
		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, domain.ApproachInfo{
			MoreMin: moreMin[i],
			RealArrivalTime: realArrivalTime[i],
			Direction: directions[i],
			ScheduledTime: scheduledTime[i],
			Delay: delay[i],
			Via: via,
			BusStop: config.GetBusStop(fetcher.from, fetcher.to),
		})
	}
	return approachInfos
}

