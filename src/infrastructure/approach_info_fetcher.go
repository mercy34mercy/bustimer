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

func (doc *CustomDocument) fetchVia(hour, min int, isHoliday bool) string {
	via := ""
	searchPath := ".timetable tr .column_day1_t2"
	if isHoliday {
		searchPath = ".timetable tr .column_day2_t2"
	}
	doc.Find(searchPath).Each(func(i int, tab *goquery.Selection) {
		if hour == (i + 5) {
			tab.Find("li").Each(func(j int, li *goquery.Selection) {
				trimed := strings.Fields(li.Text())
				if len(trimed) >= 1 {
					minute, err := strconv.Atoi(trimed[1])
					if err != nil {
						fmt.Println("Conversion failed")
					}
					if minute == min {
						via = trimed[0]
					}
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
	// 接近情報のWebから取れる情報をスクレイピングする
	customDoc := CustomDocument{approachDoc}
	moreMin, realArrivalTime, directions, scheduledTime, delay := customDoc.fetchApproachInfo()

	iterateCount := findMinLen(moreMin, realArrivalTime, directions, scheduledTime, delay)
	if len(viaUrl) > 0 {
		viaDoc, err := goquery.NewDocument(viaUrl)
		if err == nil {
			customDoc = CustomDocument{viaDoc}
		}
	}
	via := ""
	for i := 0; i < iterateCount; i++ {
		if len(viaUrl) > 0 {
			//TODO: 経由情報のスクレイピング
			//TODO: viaUrlの実装はまだ
			hour, _ := strconv.Atoi(scheduledTime[i][:2])
			min, _ := strconv.Atoi(scheduledTime[i][3:])
			via = customDoc.fetchVia(hour, min, true)
		}
		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, domain.ApproachInfo{
			MoreMin: moreMin[i],
			RealArrivalTime: realArrivalTime[i],
			Direction: directions[i],
			ScheduledTime: scheduledTime[i],
			Delay: delay[i],
			Via: via,
		})
	}
	return approachInfos
}
