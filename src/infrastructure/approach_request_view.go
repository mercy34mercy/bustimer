package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/presenter"
	"net/http"
)

func ApproachInfoRequest(c Context) error {
	// リクエストURLを作成
	approachInfoUrls, from := c.GetApproachInfoUrls()

	// URLからデータを取得するFetcherを作成
	fetcher := ApproachInfoFetcher{}

	// Presenterに処理をリクエスト
	approachInfos := presenter.RequestApproachInfos(approachInfoUrls, from, fetcher)
	return c.Response(http.StatusOK, approachInfos)
}