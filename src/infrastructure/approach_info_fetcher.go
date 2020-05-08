package infrastructure

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/domain"
	"regexp"
	"strconv"
	"strings"
)

type ApproachInfoFetcher struct {
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
					Delay = append(Delay, span.Text())
				}
			})
		})
	})
	return moreMin, realArrivalTime, Direction, ScheduledTime, Delay
}

func (doc *CustomDocument) fetchViaAndBusStop(hour, min int) (string, string) {
	via := ""
	busStop := ""
	doc.Find(".time").Each(func(i int, s *goquery.Selection) {
		h := i + 4	// HTMLの要素+3がちょうど時間(hour)になる
		if h == hour {
			// TODO: 経由とバス停情報をスクレイピングしないといけない
			// TODO: 本当は平日以外も対応したい
			s.Parent().Find(".column_day1_t2 .ttList li").Each(func(j int, s *goquery.Selection) {
				fmt.Println("s.Text(): ", s.Text())
			})
			return
		}
	})
	return via, busStop
}

func (doc *CustomDocument) fetchVia(hour, min int) string {
	via := ""
	doc.Find(".timetable tr").Each(func(tr int, s *goquery.Selection) {
		if hour == (tr + 5) {
			s.Find(".column_day1_t2 li ").Each(func(i int, li *goquery.Selection) {
				trimed := strings.TrimSpace(li.Text())
				minute, _ := strconv.Atoi(trimed)
				if minute == min {
					via = li.Find(".legend span").First().Text()
				}
			})
		}
	})
	return via
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

func (fetcher ApproachInfoFetcher) FetchApproachInfos(approachInfoUrl, viaUrl string) domain.ApproachInfos {
	approachInfos := domain.ApproachInfos{
		ApproachInfo: make([]domain.ApproachInfo, 0),
	}
	approachDoc, err := goquery.NewDocument(approachInfoUrl)
	if err != nil {
		// TODO: スクレイピングが失敗した場合の処理
		// 近江鉄道バスサーバ死亡説...
		fmt.Println(approachInfoUrl, " has no content or invalid format. unable to scrape")
		return domain.ApproachInfos{}
	}
	customDoc := CustomDocument{approachDoc}
	moreMin, realArrivalTime, directions, scheduledTime, delay := customDoc.fetchApproachInfo()

	// TODO: viaUrlの実装はまだで、とりあえずログに出してるだけ
	iterateCount := findMinLen(moreMin, realArrivalTime, directions, scheduledTime, delay)
	for i := 0; i < iterateCount; i++ {
		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, domain.ApproachInfo{
			MoreMin: moreMin[i],
			RealArrivalTime: realArrivalTime[i],
			Direction: directions[i],
			ScheduledTime: scheduledTime[i],
			Delay: delay[i],
		})
	}
	return approachInfos
}
