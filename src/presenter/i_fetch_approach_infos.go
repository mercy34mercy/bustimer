package presenter

import "github.com/shun-shun123/bus-timer/src/domain"

type IFetchApproachInfos interface {
	FetchApproachInfos(url string) domain.ApproachInfos
}
