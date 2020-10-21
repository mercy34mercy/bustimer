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
	if len(approachInfos.ApproachInfo) == 0 {
		return c.Response("ApproachInfoRequest", http.StatusNoContent, approachInfos)
	}
	return c.Response("ApproachInfoRequest", http.StatusOK, approachInfos)
}

func ApproachInfoRequestV2(c Context) error {
	approachInfoUrls := c.GetApproachInfoUrls()
	from, to := c.GetFromToQuery()

	fetcher := ApproachInfoFetcher{
		from: from,
		to: to,
	}

	approachInfos := presenter.RequestApproachInfos(approachInfoUrls, fetcher)
	v2 := approachInfos.CopyToV2()
	if len(approachInfos.ApproachInfo) == 0 {
		return c.Response("ApproachInfoRequestV2", http.StatusNoContent, v2)
	}
	return c.Response("ApproachInfoRequestV2", http.StatusOK, v2)
}