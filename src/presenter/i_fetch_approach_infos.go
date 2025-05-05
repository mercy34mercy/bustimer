package presenter

import (
	"context"

	"github.com/shun-shun123/bus-timer/src/domain"
)

type IFetchApproachInfos interface {
	FetchApproachInfos(approachInfoUrl string, pastUrlsApproachInfos domain.ApproachInfos) domain.ApproachInfos
	FetchApproachInfosWithContext(ctx context.Context, approachInfoUrl string, pastUrlsApproachInfos domain.ApproachInfos) domain.ApproachInfos
}
