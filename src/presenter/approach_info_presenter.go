package presenter

import (
	"github.com/shun-shun123/bus-timer/src/domain"
)

func RequestApproachInfos(approachInfoUrls []string, viaUrls []string, fetcher IFetchApproachInfos) domain.ApproachInfos {
	approachInfos := make([]domain.ApproachInfos, len(approachInfoUrls))
	for i, v := range approachInfoUrls {
		approachInfos[i] = fetcher.FetchApproachInfos(v, viaUrls[i])
	}
	// TODO: 上位三つの早いものを取り出す処理
	fastThree := domain.ApproachInfos{}
	// FIXME: 一旦それっぽいデータを埋めてる
	fastThree.ApproachInfo = append(fastThree.ApproachInfo, approachInfos[0].ApproachInfo[0])
	fastThree.ApproachInfo = append(fastThree.ApproachInfo, approachInfos[0].ApproachInfo[1])
	fastThree.ApproachInfo = append(fastThree.ApproachInfo, approachInfos[0].ApproachInfo[2])
	return fastThree
}
