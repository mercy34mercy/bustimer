package infrastructure

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
)

type ApproachInfoFetcher struct {
	from config.From
	to   config.To
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

	// あと約X分で発車の部分を検索し、X（残り時間）を抽出する
	doc.Find("div.text-lg.font-bold.text-error strong.mx-1.text-2xl").Each(func(i int, s *goquery.Selection) {
		waitTime := s.Text()
		moreMin = append(moreMin, waitTime)
		Direction = append(Direction, "かがやき通り")
		Delay = append(Delay, "0")
		fmt.Printf("バス %d: あと約%s分\n", i+1, waitTime)
	})

	// docを表示
	fmt.Println(doc.Document)

	// strong の tag を検索
	// doc.Find("strong").Each(func(i int, s *goquery.Selection) {
	// 	waitTime := s.Text()
	// 	moreMin = append(moreMin, waitTime)
	// 	Delay = append(Delay, "0")
	// 	fmt.Printf("バス %d: 約%s分\n", i+1, waitTime)
	// })

	arrivalTimeArray := make([]string, 0)
	doc.Find("time").Each(func(i int, s *goquery.Selection) {
		arrivaltime := s.Text()
		if strings.Contains(arrivaltime, ":") && !strings.Contains(arrivaltime, "現在") {
			arrivalTimeArray = append(arrivalTimeArray, arrivaltime)
		}
	})

	for j, v := range arrivalTimeArray {
		if j%2 == 0 {
			fmt.Println(v)
			realArrivalTime = append(realArrivalTime, v)
			ScheduledTime = append(ScheduledTime, v)
		}
	}

	// log.Print("抽出された残り時間: ", moreMin)

	// // 以下は既存の実験的なコードなので、実際の実装では不要
	// log.Print("ここまで2")
	// // var waitTimes []string
	// // <strong class="mx-1 text-2xl"> の要素を検索
	// // いくつかの異なるセレクタを試してみる
	// fmt.Println("方法1: div.text-error strong.text-2xl を使用")
	// doc.Find("div.text-error strong.text-2xl").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Printf("バス %d: あと約%s分\n", i+1, s.Text())
	// })

	// fmt.Println("\n方法2: strong.mx-1.text-2xl を使用")
	// doc.Find("strong.mx-1.text-2xl").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Printf("バス %d: あと約%s分\n", i+1, s.Text())
	// })

	// fmt.Println("\n方法3: li button div.my-3 div.font-bold strong を使用")
	// doc.Find("li button div.my-3 div.font-bold strong").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Printf("バス %d: あと約%s分\n", i+1, s.Text())
	// })
	// doc.Find(".tableDetail").Each(func(i int, s *goquery.Selection) {
	// 	if i >= dataCount {
	// 		return
	// 	}
	// 	moreMin = append(moreMin, s.Find(".more_min").Text())
	// 	r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
	// 	arrivalTime := r.FindString(s.Find(".time").Text())
	// 	realArrivalTime = append(realArrivalTime, arrivalTime)
	// 	s.Find(".bsul li").Each(func(j int, li *goquery.Selection) {
	// 		if j == 4 {
	// 			trimed := strings.TrimSpace(strings.Split(li.Text(), "\n")[2])
	// 			Direction = append(Direction, trimed)
	// 		}
	// 	})
	// 	s.Find(".moreArea .bsul li").Each(func(j int, li *goquery.Selection) {
	// 		r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
	// 		scheduledTime := r.FindStringSubmatch(li.Text())
	// 		if scheduledTime != nil {
	// 			ScheduledTime = append(ScheduledTime, scheduledTime[0])
	// 		}
	// 		li.Find(".bsmidashi").Each(func(k int, span *goquery.Selection) {
	// 			if k == 1 {
	// 				Delay = append(Delay, config.TrimParentheses(span.Text()))
	// 			}
	// 		})
	// 	})
	// })
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

func (fetcher ApproachInfoFetcher) FetchApproachInfos(approachInfoUrl string, pastUrlsApproachInfos domain.ApproachInfos) domain.ApproachInfos {
	// 返り値で返す変数を初期化
	approachInfos := domain.CreateApproachInfos()

	// var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"

	// ページの情報を取得する
	// client := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }

	// req, err := http.NewRequest("GET", approachInfoUrl, nil)
	// if err != nil {
	// 	log.Printf("Http/GET failed to %v because of %v", approachInfoUrl, err)
	// 	return approachInfos
	// }
	// // この辺りのHeaderを設定しないと403が返された
	// req.Header.Set("Accept", accept)
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	resp, err := http.Get(approachInfoUrl)
	if err != nil {
		log.Printf("Http/GET failed to %v because of %v", approachInfoUrl, err)
		return approachInfos
	}
	defer resp.Body.Close()

	log.Printf("Http/GET success to %v", approachInfoUrl)
	log.Printf("Http/GET status code: %v", resp.StatusCode)
	// レスポンスのボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("レスポンスのボディを読み込めませんでした: %v", err)
		return approachInfos
	}

	// デバッグ用にファイルに保存
	// os.WriteFile("response.html", body, 0644)

	// bodyからio.Readerを作成して2回読めるようにする
	bodyReader := strings.NewReader(string(body))

	// io.Reader経由でドキュメントにパースする
	approachDoc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader failed because of %v", err)
		return approachInfos
	}

	// html, err := approachDoc.Html()
	// log.Printf("パースされたHTML: %s", html)

	// log.Print("ここまで")

	// CustomDocument型に変換する
	customDoc := CustomDocument{approachDoc}
	moreMin, realArrivalTime, directions, scheduledTime, delay := customDoc.fetchApproachInfo()

	// どれかが空の場合もあるので、最小の数を探す
	iterateCount := findMinLen(moreMin, realArrivalTime, directions, scheduledTime, delay)
	sameTimeCountDict := map[string]int{}
	for _, pastUrlsApproachInfo := range pastUrlsApproachInfos.ApproachInfo {
		if v, ok := sameTimeCountDict[pastUrlsApproachInfo.ScheduledTime]; ok {
			sameTimeCountDict[pastUrlsApproachInfo.ScheduledTime] = v + 1
		} else {
			sameTimeCountDict[pastUrlsApproachInfo.ScheduledTime] = 0
		}
	}

	for i := 0; i < iterateCount; i++ {
		via := ""
		// hh:mmの表記でくる
		if v, ok := sameTimeCountDict[scheduledTime[i]]; ok {
			sameTimeCountDict[scheduledTime[i]] = v + 1
		} else {
			sameTimeCountDict[scheduledTime[i]] = 0
		}

		hour, _ := strconv.Atoi(scheduledTime[i][:2])
		min, _ := strconv.Atoi(scheduledTime[i][3:])
		tt, ok := TimeTableCache[fetcher.from]
		if ok {
			timeTableData := tt.Weekdays
			if config.IsHoliday() {
				timeTableData = tt.Saturdays
			}
			for j, v := range timeTableData[hour] {
				if convMin, err := strconv.Atoi(v.Min); err == nil {
					if convMin == min {
						if j+sameTimeCountDict[scheduledTime[i]] >= len(timeTableData[hour]) {
							break
						}
						var matchedOneBusTime = timeTableData[hour][j+sameTimeCountDict[scheduledTime[i]]]
						matchMin, _ := strconv.Atoi(matchedOneBusTime.Min)
						if matchMin == min {
							via = timeTableData[hour][j+sameTimeCountDict[scheduledTime[i]]].Via
						}
						break
					}
				}
			}
		} else if fetcher.from != config.Unknown {
			via = config.GetVia(fetcher.from)
		}
		//FIXME: 特例の処理
		if fetcher.from == config.FromMinakusa && scheduledTime[i] == "17:20" {
			continue
		}
		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, domain.ApproachInfo{
			MoreMin:         moreMin[i],
			RealArrivalTime: realArrivalTime[i],
			Direction:       directions[i],
			ScheduledTime:   scheduledTime[i],
			Delay:           delay[i],
			Via:             via,
			BusStop:         config.GetApproachBusStop(fetcher.from, fetcher.to, via),
			RequiredTime:    config.GetRequiredTime(fetcher.from, fetcher.to, via),
		})
	}
	return approachInfos
}
