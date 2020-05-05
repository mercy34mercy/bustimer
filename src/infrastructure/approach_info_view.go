package infrastructure

import (
	"github.com/labstack/echo/v4"
	"github.com/shun-shun123/bus-timer/src/domain"
	"github.com/shun-shun123/bus-timer/src/presenter"
	"net/http"
)

func fetchApproachInfo() (domain.ApproachInfos, error) {
	html := ""
	approachInfos, err := presenter.HtmlToApproach(html)
	if err != nil {
		return approachInfos, err
	}
	return approachInfos, nil
}


func RequestApproachInfo(c echo.Context) error {
	approachInfo, err := fetchApproachInfo()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, approachInfo)
}