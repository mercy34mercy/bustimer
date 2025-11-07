package infrastructure

import (
	"context"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ApproachInfoFetcher struct {
	from config.From
	to   config.To
}

type CustomDocument struct {
	*goquery.Document
}

var re = regexp.MustCompile(`約([0-9]+)分`)
var fetcherTracer = otel.Tracer("bustimer/fetcher")

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
	return fetcher.FetchApproachInfosWithContext(context.Background(), approachInfoUrl, pastUrlsApproachInfos)
}

func (fetcher ApproachInfoFetcher) FetchApproachInfosWithContext(ctx context.Context, approachInfoUrl string, pastUrlsApproachInfos domain.ApproachInfos) domain.ApproachInfos {
	// トレーシング
	ctx, span := fetcherTracer.Start(ctx, "FetchApproachInfosWithContext")
	defer span.End()

	startTime := time.Now()

	span.SetAttributes(
		attribute.String("url", approachInfoUrl),
		attribute.String("from", fetcher.from.String()),
		attribute.String("to", fetcher.to.String()),
	)

	// 返り値で返す変数を初期化
	approachInfos := domain.CreateApproachInfos()

	// HTTPリクエスト
	_, httpSpan := fetcherTracer.Start(ctx, "HTTPRequest")

	// HTTPクライアントを作成してヘッダーを設定
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", approachInfoUrl, nil)
	if err != nil {
		log.Printf("Failed to create request for %v: %v", approachInfoUrl, err)
		httpSpan.SetAttributes(attribute.String("error", err.Error()))
		httpSpan.End()
		return approachInfos
	}

	// ブラウザを模倣するヘッダーを追加
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://transfer-cloud.navitime.biz/")

	resp, err := client.Do(req)
	httpSpan.SetAttributes(attribute.String("url", approachInfoUrl))

	if err != nil {
		log.Printf("Http/GET failed to %v because of %v", approachInfoUrl, err)
		httpSpan.SetAttributes(attribute.String("error", err.Error()))
		httpSpan.End()
		return approachInfos
	}
	defer resp.Body.Close()

	// HTTPステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		// レスポンスボディを読み込んでログに記録
		body, _ := io.ReadAll(resp.Body)
		bodyPreview := string(body)
		if len(bodyPreview) > 500 {
			bodyPreview = bodyPreview[:500] + "..."
		}
		log.Printf("HTTP request to %v returned status code %d. Response body: %s", approachInfoUrl, resp.StatusCode, bodyPreview)
		httpSpan.SetAttributes(
			attribute.Int("statusCode", resp.StatusCode),
			attribute.String("status", resp.Status),
		)
		httpSpan.End()
		return approachInfos
	}
	httpSpan.SetAttributes(attribute.Int("statusCode", resp.StatusCode))
	httpSpan.End()

	// レスポンスのボディを読み込む
	_, parseSpan := fetcherTracer.Start(ctx, "ParseResponse")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("レスポンスのボディを読み込めませんでした: %v", err)
		parseSpan.SetAttributes(attribute.String("error", err.Error()))
		parseSpan.End()
		return approachInfos
	}

	// bodyからio.Readerを作成して2回読めるようにする
	bodyReader := strings.NewReader(string(body))

	// io.Reader経由でドキュメントにパースする
	approachDoc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader failed because of %v", err)
		parseSpan.SetAttributes(attribute.String("error", err.Error()))
		parseSpan.End()
		return approachInfos
	}
	parseSpan.End()

	// ドキュメントを操作するカスタムラッパー
	doc := &CustomDocument{approachDoc}

	// スクレイピング
	_, scrapeSpan := fetcherTracer.Start(ctx, "Scraping")
	moreMin, realArrivalTime, direction, scheduledTime, delay, busstop, via, requiredTime := doc.FetchApproachInfo()

	// 結果の件数をトレース
	scrapeSpan.SetAttributes(
		attribute.Int("moreMinCount", len(moreMin)),
		attribute.Int("realArrivalTimeCount", len(realArrivalTime)),
		attribute.Int("directionCount", len(direction)),
		attribute.Int("scheduledTimeCount", len(scheduledTime)),
		attribute.Int("delayCount", len(delay)),
		attribute.Int("busstopCount", len(busstop)),
		attribute.Int("viaCount", len(via)),
		attribute.Int("requiredTimeCount", len(requiredTime)),
	)
	scrapeSpan.End()

	// 最小の長さを取得 (viaは除外 - オプショナルなフィールドのため)
	min := findMinLenWithIntSlice(requiredTime, moreMin, realArrivalTime, direction, scheduledTime, delay, busstop)

	// スクレイピング結果が空の場合、警告ログを出力
	if min == 0 {
		log.Printf("Warning: No bus approach info found from %v. Scraped data counts - moreMin:%d, realArrivalTime:%d, direction:%d, scheduledTime:%d, delay:%d, busstop:%d, via:%d, requiredTime:%d, calculated min:%d",
			approachInfoUrl, len(moreMin), len(realArrivalTime), len(direction), len(scheduledTime), len(delay), len(busstop), len(via), len(requiredTime), min)

		// no-contentsセクションをチェック
		noContentFound := false
		noContentMessage := ""
		doc.Find(".no-contents").Each(func(i int, s *goquery.Selection) {
			noContentFound = true
			noContentMessage = s.Text()
		})

		if noContentFound {
			log.Printf("Info: 'no-contents' section detected. Message: %s", noContentMessage)
		}

		// デバッグ: HTMLレスポンスの一部を出力
		bodyPreview := string(body)
		if len(bodyPreview) > 1000 {
			bodyPreview = bodyPreview[:1000]
		}
		log.Printf("Debug: HTML response preview (first 1000 chars):\n%s", bodyPreview)

		// デバッグ: 主要なセレクタの結果を確認
		log.Printf("Debug: Checking selectors - div.text-lg.font-bold.text-error count: %d", doc.Find("div.text-lg.font-bold.text-error").Length())
		log.Printf("Debug: Checking selectors - time count: %d", doc.Find("time").Length())
		log.Printf("Debug: Checking selectors - button.w-full.rounded count: %d", doc.Find("button.w-full.rounded.text-left.drop-shadow-md").Length())
	} else {
		log.Printf("Info: Successfully scraped %d bus approach info(s) from %v. Scraped data counts - moreMin:%d, realArrivalTime:%d, direction:%d, scheduledTime:%d, delay:%d, busstop:%d, via:%d, requiredTime:%d",
			min, approachInfoUrl, len(moreMin), len(realArrivalTime), len(direction), len(scheduledTime), len(delay), len(busstop), len(via), len(requiredTime))
	}

	for i := 0; i < min; i++ {
		// viaが範囲外の場合は空文字列を使用
		viaValue := ""
		if i < len(via) {
			viaValue = via[i]
		}

		info := domain.ApproachInfo{
			MoreMin:         moreMin[i],
			RealArrivalTime: realArrivalTime[i],
			Direction:       direction[i],
			Via:             viaValue,
			ScheduledTime:   scheduledTime[i],
			Delay:           delay[i],
			BusStop:         busstop[i],
			RequiredTime:    requiredTime[i],
		}
		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, info)
	}

	// 処理時間の記録
	processingTime := time.Since(startTime)
	span.SetAttributes(
		attribute.String("processingTime", processingTime.String()),
		attribute.Int64("processingTimeMs", processingTime.Milliseconds()),
		attribute.Int("resultCount", len(approachInfos.ApproachInfo)),
	)

	return approachInfos
}
