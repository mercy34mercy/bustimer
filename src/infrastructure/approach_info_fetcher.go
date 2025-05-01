package infrastructure

import (
	"fmt"
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

var dataCount = 3

// 接近情報のサイトから取れる下記の情報をまとめてスクレイピングする
func (doc *CustomDocument) fetchApproachInfo() ([]string, []string, []string, []string, []string, []string, []string, []int) {
	moreMin := make([]string, 0)
	realArrivalTime := make([]string, 0)
	Direction := make([]string, 0)
	ScheduledTime := make([]string, 0)
	Busstop := make([]string, 0)
	Via := make([]string, 0)
	RequiredTime := make([]int, 0)

	doc.Find("div.text-sm.mb-2.ml-auto.mr-4.w-fit.text-text-grey").Each(func(i int, s *goquery.Selection) {
		requiredTime := s.Text()
		// 予測所要時間 約20分 こういうテキストだから、数字の部分を抽出する
		re := regexp.MustCompile(`約([0-9]+)分`)
		matched := re.FindStringSubmatch(requiredTime)
		if len(matched) > 1 {
			requiredTimeInt, err := strconv.Atoi(matched[1])
			if err != nil {
				RequiredTime = append(RequiredTime, 0)
			} else {
				RequiredTime = append(RequiredTime, requiredTimeInt)
			}
		} else {
			RequiredTime = append(RequiredTime, 0)
		}
	})

	// あと約X分で発車の部分を検索し、X（残り時間）を抽出する
	doc.Find("div.text-lg.font-bold.text-error strong.mx-1.text-2xl").Each(func(i int, s *goquery.Selection) {
		waitTime := s.Text()
		moreMin = append(moreMin, waitTime)
	})

	doc.Find("dt.mr-1.break-all").Each(func(i int, s *goquery.Selection) {
		busstop := s.Text()
		// 最後の1文字を抽出
		lastChar := busstop[len(busstop)-1:]
		Busstop = append(Busstop, lastChar)
	})
	Delay := make([]string, len(moreMin))
	doc.Find("div.flex.justify-center.text-error").Each(func(i int, s *goquery.Selection) {
		println(s.Text())
		s.Find("span.mr-2").Each(func(i int, s *goquery.Selection) {
			// textの2文字目を抽出
			if len(s.Text()) > 1 {
				delay := s.Text()[1]
				Delay = append(Delay, string(delay))
			} else {
				Delay = append(Delay, "")
			}
		})
	})
	doc.Find("button.w-full.rounded.text-left.drop-shadow-md.bg-white").Each(func(i int, s *goquery.Selection) {
		s.Find("div.flex.justify-between").Each(func(i int, t *goquery.Selection) {
			t.Find("div.flex.flex-col").Each(func(i int, s *goquery.Selection) {
				s.Find("span").Each(func(i int, s *goquery.Selection) {
					fmt.Println("via: ", s.Text())
					via := s.Text()
					Via = append(Via, via)
				})
			})
		})
	})
	doc.Find("div.flex.flex-col").Each(func(i int, s *goquery.Selection) {
		s.Find("span.font-bold").Each(func(i int, s *goquery.Selection) {
			dir := s.Text()
			Direction = append(Direction, dir)
			// Directionの[]で囲まれた文字列を抽出
			fmt.Println(dir)
			// // [ or ]でsplitして、[ or ]の中身を抽出
			// via := strings.Split(dir, "[")
			// if len(via) > 1 {
			// 	// viaの最後の1文字以外を抽出
			// 	Via = append(Via, via[1][:len(via[1])-1])
			// } else {
			// 	Via = append(Via, "")
			// }
		})
	})

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
	return moreMin, realArrivalTime, Direction, ScheduledTime, Delay, Busstop, Via, RequiredTime
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
	moreMin, realArrivalTime, directions, scheduledTime, delay, busstop, via, requiredTime := customDoc.fetchApproachInfo()

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
			Via:             via[i],
			BusStop:         busstop[i],
			RequiredTime:    requiredTime[i],
		})
	}
	return approachInfos
}
