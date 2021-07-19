package presenter

import (
	"github.com/shun-shun123/bus-timer/src/domain"
)

func RequestApproachInfos(approachInfoUrls []string, fetcher IFetchApproachInfos) domain.ApproachInfos {
	approachInfos := domain.CreateApproachInfos()
	for _, url := range approachInfoUrls {
		fetchResult := fetcher.FetchApproachInfos(url, approachInfos)
		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, fetchResult.ApproachInfo...)
	}
	fastThree := approachInfos.GetFastThree()
	return fastThree
}
