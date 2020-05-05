package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/presenter"
	"net/http"
)

func ApproachInfoRequest(c Context) error {
	urls := c.GetApproachInfoUrl()
	fetcher := ApproachInfoFetcher{}
	approachInfos := presenter.RequestApproachInfos(urls, fetcher)
	return c.Response(http.StatusOK, approachInfos)
}