package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/presenter"
	"net/http"
)

func ApproachInfoRequest(c Context) error {
	// リクエストURLを作成
	approachInfoUrls := c.GetApproachInfoUrls()
	from, to := c.GetFromToQuery()

	// URLからデータを取得するFetcherを作成
	fetcher := ApproachInfoFetcher{
		from: from,
		to: to,
	}

	// Presenterに処理をリクエスト
	approachInfos := presenter.RequestApproachInfos(approachInfoUrls, fetcher)
	return c.Response(http.StatusOK, approachInfos)
}