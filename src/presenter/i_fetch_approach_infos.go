package presenter

import (
	"github.com/shun-shun123/bus-timer/src/domain"
)

type IFetchApproachInfos interface {
	FetchApproachInfos(approachInfoUrl string) domain.ApproachInfos
}
