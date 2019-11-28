package approach

import (
	"fmt"
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
	if fr == frRits {
		for i := 0; i < MAX_RESPONSE; i++ {
			fullURL := url + "?fr=" + fr + "&dgmpl=" + dgmpl[0]
			info, err := scrapeFromURL(fullURL, i)
			if err == nil && info.MoreMin != "" {
				approach["res"] = append(approach["res"], info)
			}
			c.Echo().Logger.Debug("Scrape from " + fullURL)
		}
	} else if fr == frMinakusa {
		for _, v := range dgmpl {
			fullURL := url + "?fr=" + fr + "&dgmpl=" + v
			if v == dgmplMap[frMinakusa][2] && len(approach["res"]) != MAX_RESPONSE-1 {
				for i := 0; i < MAX_RESPONSE-len(approach["res"]); i++ {
					info, err := scrapeFromURL(fullURL, i)
					if err == nil && info.MoreMin != "" {
						approach["res"] = append(approach["res"], info)
					}
					c.Echo().Logger.Debug("Scrape from " + fullURL)
				}
			} else {
				info, err := scrapeFromURL(fullURL, 0)
				if err == nil && info.MoreMin != "" {
					approach["res"] = append(approach["res"], info)
				}
				c.Echo().Logger.Debug("Scrape From " + fullURL)
			}
		}
	}
	// 接近情報があるかどうかを判断する
	if len(approach["res"]) == 0 {
		return c.JSON(http.StatusNoContent, approach)
	}
	return c.JSON(http.StatusOK, approach)
}

func scrapeFromURL(fullURL string, index int) (approachInfo, error) {
	approachInfo := approachInfo{}
	doc, err := goquery.NewDocument(fullURL)
	if err != nil {
		return approachInfo, err
	}
	approachInfo, err = scrapeApproachInfo(doc, index)
	if err != nil {
		return approachInfo, err
	}
	return approachInfo, nil
}

func scrapeApproachInfo(doc *goquery.Document, index int) (approachInfo, error) {
	approachInfo := approachInfo{}
	doc.Find(".more_min").Each(func(i int, s *goquery.Selection) {
		if i == index {
			target := strings.TrimSpace(s.Text())
			approachInfo.MoreMin = target
			fmt.Println(target)
		}
	})
	doc.Find(".time").Each(func(i int, s *goquery.Selection) {
		if i == index {
			r := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
			target := r.FindStringSubmatch(strings.TrimSpace(s.Text()))
			approachInfo.RealArrivalTime = target[0]
			fmt.Println(target[0])
		}
	})
	doc.Find(".tableDetail").Each(func(i int, s *goquery.Selection) {
		if i == index {
			s.Find(".bsul").First().Find("li").Each(func(j int, li *goquery.Selection) {
				switch j {
				case 3:
					trimed := strings.Trim(strings.Trim(li.Text(), "\n"), " ")
					splited := strings.Split(trimed, "\n")
					content := strings.Trim(splited[1], " ")
					approachInfo.Descent = content
				case 4:
					trimed := strings.Trim(strings.Trim(li.Text(), "\n"), " ")
					splited := strings.Split(trimed, "\n")
					content := strings.Trim(splited[1], " ")
					approachInfo.Direction = content
				}
			})
		}
	})
	return approachInfo, nil
}
