package approach

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
)

const (
	url = "https://ohmitetudo-bus.jorudan.biz/busstatedtl"
)

var dgmplMap = map[string][]string{"南草津駅〔近江鉄道・湖国バス〕": []string{"南草津駅〔近江鉄道・湖国バス〕:1", "南草津駅〔近江鉄道・湖国バス〕:3", "南草津駅〔近江鉄道・湖国バス〕:4"},
	"立命館大学〔近江鉄道・湖国バス〕": []string{"立命館大学〔近江鉄道・湖国バス〕:2"}}

type approachInfo struct {
	// あと何分で到着か
	MoreMin string `json:"more_min"`
	// 実際の到着予定時刻(遅延を考慮する)
	RealArrivalTime string `json:"real_arrive_time"`
	// 系統
	Descent string `json:"descent"`
	// 行き先
	Direction string `json:"direction"`
}

// ScrapeApproachInfo バスの接近情報をスクレイピングする
func ScrapeApproachInfo(c echo.Context) error {
	fr := c.QueryParam("fr")
	dgmpl := dgmplMap[fr]
	approachInfo := map[string][]approachInfo{}
	for _, v := range dgmpl {
		fullURL := url + "?fr=" + fr + "&dgmpl=" + v
		info, err := scrapeFromURL(fullURL)
		if err == nil {
			approachInfo["res"] = append(approachInfo["res"], info)
		}
		c.Echo().Logger.Debug("Scrape from " + fullURL)
	}
	return c.JSON(http.StatusOK, approachInfo)
}

func scrapeFromURL(fullURL string) (approachInfo, error) {
	approachInfo := approachInfo{}
	doc, err := goquery.NewDocument(fullURL)
	if err != nil {
		return approachInfo, err
	}
	approachInfo, err = scrapeApproachInfo(doc)
	if err != nil {
		return approachInfo, err
	}
	return approachInfo, nil
}

func scrapeApproachInfo(doc *goquery.Document) (approachInfo, error) {
	approachInfo := approachInfo{}
	doc.Find(".bsul").First().Find("li").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 1:
			approachInfo.MoreMin = strings.TrimSpace(s.Text())
		case 2:
			approachInfo.RealArrivalTime = strings.TrimSpace(s.Text())
		case 3:
			trimed := strings.Trim(strings.Trim(s.Text(), "\n"), " ")
			splited := strings.Split(trimed, "\n")
			approachInfo.Descent = strings.Trim(splited[1], " ")
		case 4:
			trimed := strings.Trim(strings.Trim(s.Text(), "\n"), " ")
			splited := strings.Split(trimed, "\n")
			approachInfo.Direction = strings.Trim(splited[1], " ")
		}
	})
	return approachInfo, nil
}
