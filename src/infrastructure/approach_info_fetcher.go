package infrastructure

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/domain"
	"regexp"
	"strings"
)

type ApproachInfoFetcher struct {
}

type CustomDocument struct {
	*goquery.Document
}

var dataCount = 3

func (doc *CustomDocument) fetchMoreMin() []string {
	moreMin := make([]string, 0)
	doc.Find(".more_min").Each(func(i int, s *goquery.Selection) {
		if i >= dataCount {
			return
		}
		target := strings.TrimSpace(s.Text())
		moreMin = append(moreMin, target)
	})
	return moreMin
}

func (doc *CustomDocument) fetchRealArrivalTime() []string {
	realArrivalTime := make([]string, 0)
	doc.Find(".time").Each(func(i int, s *goquery.Selection) {
		if i >= dataCount {
			return
		}
		r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
		target := r.FindStringSubmatch(strings.TrimSpace(s.Text()))
		realArrivalTime = append(realArrivalTime, target[0])
	})
	return realArrivalTime
}

func (doc *CustomDocument) fetchDirection() []string {
	directions := make([]string, 0)
	doc.Find(".tableDetail").Each(func(i int, s *goquery.Selection) {
		if i >= dataCount {
			return
		}
		s.Find(".bsul").First().Find("li").Each(func(j int, li *goquery.Selection) {
			if j == 4 {
				trimed := strings.Trim(strings.Trim(li.Text(), "\n"), " ")
				splited := strings.Split(trimed, "\n")
				content := strings.Trim(splited[1], " ")
				directions = append(directions, content)
			}
		})
	})
	return directions
}

func (doc *CustomDocument) fetchScheduledTimeAndDelay() ([]string, []string) {
	scheduledTime := make([]string, 0)
	delay := make([]string, 0)
	doc.Find(".moreArea").Each(func(i int, s *goquery.Selection) {
		if i >= dataCount {
			return
		}
		s.Find(".bsmidashi").Each(func(j int, li *goquery.Selection) {
			if j == 0 {
				r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
				target := r.FindStringSubmatch(li.Parent().Text())
				scheduledTime = append(scheduledTime, target[0])
			} else if j == 1 {
				delay = append(delay, li.Text())
			}
		})
	})
	return scheduledTime, delay
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

func (fetcher ApproachInfoFetcher) FetchApproachInfos(approachInfoUrl, viaUrl string) domain.ApproachInfos {
	approachInfos := domain.ApproachInfos{
		ApproachInfo: make([]domain.ApproachInfo, 3),
	}
	approachDoc, err := goquery.NewDocument(approachInfoUrl)
	if err != nil {
		// TODO: スクレイピングが失敗した場合の処理
		// 近江鉄道バスサーバ死亡説...
		log.Fatal(approachInfoUrl, " has no content or invalid format. unable to scrape")
		return domain.ApproachInfos{}
	}
	customDoc := CustomDocument{approachDoc}
	moreMin := customDoc.fetchMoreMin()
	realArrivalTime := customDoc.fetchRealArrivalTime()
	directions := customDoc.fetchDirection()
	scheduledTime, delay := customDoc.fetchScheduledTimeAndDelay()

	// TODO: viaUrlの実装はまだで、とりあえずログに出してるだけ
	for i := 0; i < dataCount; i++ {
		approachInfos.ApproachInfo[i].MoreMin = moreMin[i]
		approachInfos.ApproachInfo[i].RealArrivalTime = realArrivalTime[i]
		approachInfos.ApproachInfo[i].Direction = directions[i]
		approachInfos.ApproachInfo[i].ScheduledTime = scheduledTime[i]
		approachInfos.ApproachInfo[i].Delay = delay[i]
	}
	return approachInfos
}
