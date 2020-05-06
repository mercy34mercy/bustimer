package infrastructure

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/domain"
	"regexp"
	"strings"
)

type ApproachInfoFetcher struct {
}

func (fetcher ApproachInfoFetcher) FetchApproachInfos(url string) domain.ApproachInfos {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		// TODO: スクレイピングが失敗した場合の処理
		// 近江鉄道バスサーバ死亡説...
		return domain.ApproachInfos{}
	}
	fmt.Println(url)
	fmt.Println(doc)
	approachInfos := domain.ApproachInfos{
		ApproachInfo: make([]domain.ApproachInfo, 3),
	}
	// さぁ、スクレイピングと行こうか。
	// TODO: スクレイピングが成功したら情報を取り出す
	doc.Find(".more_min").Each(func(i int, s *goquery.Selection) {
		if i >= 3 {
			return
		}
		target := strings.TrimSpace(s.Text())
		approachInfos.ApproachInfo[i].MoreMin = target
	})
	doc.Find(".time").Each(func(i int, s *goquery.Selection) {
		if i >= 3 {
			return
		}
		r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
		target := r.FindStringSubmatch(strings.TrimSpace(s.Text()))
		approachInfos.ApproachInfo[i].RealArrivalTime = target[0]
	})
	doc.Find(".tableDetail").Each(func(i int, s *goquery.Selection) {
		if i >= 3 {
			return
		}
		s.Find(".bsul").First().Find("li").Each(func(j int, li *goquery.Selection) {
			if j == 4 {
				trimed := strings.Trim(strings.Trim(li.Text(), "\n"), " ")
				splited := strings.Split(trimed, "\n")
				content := strings.Trim(splited[1], " ")
				approachInfos.ApproachInfo[i].Direction = content
			}
		})
	})

	doc.Find(".moreArea").Each(func(i int, s *goquery.Selection) {
		if i >= 3 {
			return
		}
		s.Find(".bsmidashi").Each(func(j int, li *goquery.Selection) {
			if j == 0 {
				r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
				target := r.FindStringSubmatch(li.Parent().Text())
				approachInfos.ApproachInfo[i].ScheduledTime = target[0]
			} else if j == 1 {
				approachInfos.ApproachInfo[i].Delay = li.Text()
			}
		})
	})
	return approachInfos
}
