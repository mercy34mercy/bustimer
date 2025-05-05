package infrastructure

import (
	"io"
	"log"
	"net/http"
	"regexp"
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

var re = regexp.MustCompile(`約([0-9]+)分`)

// FetchApproachInfo 接近情報のサイトから取れる下記の情報をまとめてスクレイピングする
func (doc *CustomDocument) FetchApproachInfo() ([]string, []string, []string, []string, []string, []string, []string, []int) {
	moreMin := doc.FetchMoreMin()
	realArrivalTime, scheduledTime := doc.FetchArrivalTime()
	direction := doc.FetchDirection()
	busstop := doc.FetchBusStop()
	delay := doc.FetchDelay(len(moreMin))
	via := doc.FetchVia()
	requiredTime := doc.FetchRequiredTime()

	return moreMin, realArrivalTime, direction, scheduledTime, delay, busstop, via, requiredTime
}

// FetchRequiredTime 所要時間の取得
func (doc *CustomDocument) FetchRequiredTime() []int {
	requiredTime := make([]int, 0)
	doc.Find("div.text-sm.mb-2.ml-auto.mr-4.w-fit.text-text-grey").Each(func(i int, s *goquery.Selection) {
		requiredTimeText := s.Text()
		// 予測所要時間 約20分 こういうテキストだから、数字の部分を抽出する
		matched := re.FindStringSubmatch(requiredTimeText)
		if len(matched) > 1 {
			requiredTimeInt, err := strconv.Atoi(matched[1])
			if err != nil {
				requiredTime = append(requiredTime, 0)
			} else {
				requiredTime = append(requiredTime, requiredTimeInt)
			}
		} else {
			requiredTime = append(requiredTime, 0)
		}
	})
	return requiredTime
}

// FetchMoreMin あと約X分で発車の部分を検索し、X（残り時間）を抽出する
func (doc *CustomDocument) FetchMoreMin() []string {
	moreMin := make([]string, 0)
	doc.Find("div.text-lg.font-bold.text-error strong.mx-1.text-2xl").Each(func(i int, s *goquery.Selection) {
		waitTime := s.Text()
		moreMin = append(moreMin, waitTime)
	})
	return moreMin
}

// FetchBusStop バス停の取得
func (doc *CustomDocument) FetchBusStop() []string {
	busstop := make([]string, 0)
	doc.Find("dt.mr-1.break-all").Each(func(i int, s *goquery.Selection) {
		busstopText := s.Text()
		// 最後の1文字を抽出
		lastChar := busstopText[len(busstopText)-1:]
		busstop = append(busstop, lastChar)
	})
	return busstop
}

// FetchDelay 遅延情報の取得
func (doc *CustomDocument) FetchDelay(baseLen int) []string {
	delay := make([]string, baseLen)
	doc.Find("div.flex.justify-center.text-error").Each(func(i int, s *goquery.Selection) {
		// println(s.Text()) - コメントアウト
		s.Find("span.mr-2").Each(func(i int, s *goquery.Selection) {
			// textの2文字目を抽出
			if len(s.Text()) > 1 {
				delayChar := s.Text()[1]
				delay = append(delay, string(delayChar))
			} else {
				delay = append(delay, "")
			}
		})
	})
	return delay
}

// FetchVia 経由地の取得
func (doc *CustomDocument) FetchVia() []string {
	via := make([]string, 0)
	doc.Find("button.w-full.rounded.text-left.drop-shadow-md.bg-white").Each(func(i int, s *goquery.Selection) {
		s.Find("div.flex.justify-between").Each(func(i int, t *goquery.Selection) {
			t.Find("div.flex.flex-col").Each(func(i int, s *goquery.Selection) {
				s.Find("span").Each(func(i int, s *goquery.Selection) {
					// fmt.Println("via: ", s.Text()) - コメントアウト
					viaText := s.Text()
					via = append(via, viaText)
				})
			})
		})
	})
	return via
}

// FetchDirection 方向の取得
func (doc *CustomDocument) FetchDirection() []string {
	direction := make([]string, 0)
	doc.Find("div.flex.flex-col").Each(func(i int, s *goquery.Selection) {
		s.Find("span.font-bold").Each(func(i int, s *goquery.Selection) {
			dir := s.Text()
			direction = append(direction, dir)
			// fmt.Println(dir) - コメントアウト
		})
	})
	return direction
}

// FetchArrivalTime 到着時間の取得
func (doc *CustomDocument) FetchArrivalTime() ([]string, []string) {
	realArrivalTime := make([]string, 0)
	scheduledTime := make([]string, 0)
	arrivalTimeArray := make([]string, 0)

	doc.Find("time").Each(func(i int, s *goquery.Selection) {
		arrivaltime := s.Text()
		if strings.Contains(arrivaltime, ":") && !strings.Contains(arrivaltime, "現在") {
			arrivalTimeArray = append(arrivalTimeArray, arrivaltime)
		}
	})

	for j, v := range arrivalTimeArray {
		if j%2 == 0 {
			// fmt.Println(v) - コメントアウト
			realArrivalTime = append(realArrivalTime, v)
			scheduledTime = append(scheduledTime, v)
		}
	}

	return realArrivalTime, scheduledTime
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

// 整数スライスの長さを考慮する関数を追加
func findMinLenWithIntSlice(intSlice []int, dataset ...[]string) int {
	min := findMinLen(dataset...)
	if len(intSlice) < min {
		min = len(intSlice)
	}
	return min
}

func (fetcher ApproachInfoFetcher) FetchApproachInfos(approachInfoUrl string, pastUrlsApproachInfos domain.ApproachInfos) domain.ApproachInfos {
	// 返り値で返す変数を初期化
	approachInfos := domain.CreateApproachInfos()

	resp, err := http.Get(approachInfoUrl)
	if err != nil {
		log.Printf("Http/GET failed to %v because of %v", approachInfoUrl, err)
		return approachInfos
	}
	defer resp.Body.Close()

	// レスポンスのボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("レスポンスのボディを読み込めませんでした: %v", err)
		return approachInfos
	}

	// bodyからio.Readerを作成して2回読めるようにする
	bodyReader := strings.NewReader(string(body))

	// io.Reader経由でドキュメントにパースする
	approachDoc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader failed because of %v", err)
		return approachInfos
	}

	// CustomDocument型に変換する
	customDoc := CustomDocument{approachDoc}
	moreMin, realArrivalTime, directions, scheduledTime, delay, busstop, via, requiredTime := customDoc.FetchApproachInfo()

	// どれかが空の場合もあるので、最小の数を探す
	iterateCount := findMinLenWithIntSlice(requiredTime, moreMin, realArrivalTime, directions, scheduledTime, delay, via, busstop)
	sameTimeCountDict := map[string]int{}
	for _, pastUrlsApproachInfo := range pastUrlsApproachInfos.ApproachInfo {
		if v, ok := sameTimeCountDict[pastUrlsApproachInfo.ScheduledTime]; ok {
			sameTimeCountDict[pastUrlsApproachInfo.ScheduledTime] = v + 1
		} else {
			sameTimeCountDict[pastUrlsApproachInfo.ScheduledTime] = 0
		}
	}

	for i := 0; i < iterateCount; i++ {
		// hh:mmの表記でくる
		if v, ok := sameTimeCountDict[scheduledTime[i]]; ok {
			sameTimeCountDict[scheduledTime[i]] = v + 1
		} else {
			sameTimeCountDict[scheduledTime[i]] = 0
		}
		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, domain.ApproachInfo{
			MoreMin:         moreMin[i],
			RealArrivalTime: realArrivalTime[i],
			Direction:       directions[i],
			ScheduledTime:   scheduledTime[i],
			Delay:           delay[i],
			Via:             via[i],
			BusStop:         busstop[i],
			RequiredTime:    requiredTime[i],
		})
	}
	return approachInfos
}
