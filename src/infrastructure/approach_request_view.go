package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/presenter"
	"net/http"
)

func ApproachInfoRequest(c Context) error {
	// リクエストURLを作成
	urls := c.GetApproachInfoUrl()

	// URLからデータを取得するFetcherを作成
	fetcher := ApproachInfoFetcher{}

	// Presenterに処理をリクエスト
	approachInfos := presenter.RequestApproachInfos(urls, fetcher)
	return c.Response(http.StatusOK, approachInfos)
}