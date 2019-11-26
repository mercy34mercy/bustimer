package approach

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

// ScrapeApproachInfo バスの接近情報をスクレイピングする
func ScrapeApproachInfo(c echo.Context) error {
	fr := c.QueryParam("fr")
	dgmpl := dgmplMap[fr]
	if len(dgmpl) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"result": "frクエリに対応するdgmplが定義されていません。frクエリの内容を確認してください"})
	}
	approach := map[string][]approachInfo{}
	for _, v := range dgmpl {
		fullURL := url + "?fr=" + fr + "&dgmpl=" + v
		info, err := scrapeFromURL(fullURL)
		if err == nil && info.MoreMin != "" {
			approach["res"] = append(approach["res"], info)
		}
		c.Echo().Logger.Debug("Scrape from " + fullURL)
	}
	// 接近情報があるかどうかを判断する
	if len(approach["res"]) == 0 {
		return c.JSON(http.StatusNoContent, approach)
	}
	return c.JSON(http.StatusOK, approach)
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
			r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
			approachInfo.RealArrivalTime = r.FindStringSubmatch(strings.TrimSpace(s.Text()))[0]
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
