package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/slack"
	"net/http"
)

type CustomContext struct {
	echo.Context
}

// echo.Contextからクエリ情報を取り出して接近情報のスクレイピングをするためのURLを生成する
func (c *CustomContext) GetApproachInfoUrls() []string {
	// クエリの抽出
	fr := c.Context.QueryParam("fr")
	to := c.Context.QueryParam("to")

	// リクエストクエリからスクレイピング用のURLに含めるクエリに変換する
	from := config.GetFromKey(fr)
	dgmplList := config.GetDgmplList(fr, to)

	// URLスライス作成（複数ある場合があるので）
	approachInfoUrls := make([]string, 0)

	// URLを作成
	for _, dgmpl := range dgmplList {
		approachInfoUrls = append(approachInfoUrls, config.CreateApproachInfoUrl(from, dgmpl))
	}
	return approachInfoUrls
}

// echo.Context経由のレスポンスをラップした
func (c *CustomContext) Response(methodName string, statusCode int, param interface{}) error {
	errorNotification(statusCode, methodName)
	return c.JSON(statusCode, param)
}

func errorNotification(statusCode int, methodName string) {
	switch (statusCode) {
	case http.StatusBadRequest:
		slack.PostMessage(fmt.Sprintf("%s StatusBadRequest", methodName))
	case http.StatusNoContent:
		slack.PostMessage(fmt.Sprintf("%s StatusNoContent", methodName))
	}
}

// 「どこ発のどこ行き」かを判定する
func (c *CustomContext) GetFromToQuery() (config.From, config.To) {
	fr := c.Context.QueryParam("fr")
	toQuery := c.Context.QueryParam("to")

	from := config.ConvertFromFrQuery(fr)
	to := config.ConvertFromToQuery(toQuery)
	return from, to
}