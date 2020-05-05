package presenter

import (
	"github.com/shun-shun123/bus-timer/src/domain"
)

func RequestApproachInfos(urls []string, fetcher IFetchApproachInfos) domain.ApproachInfos {
	approachInfos := make([]domain.ApproachInfos, len(urls))
	for i, v := range urls {
		approachInfos[i] = fetcher.FetchApproachInfos(v)
	}
	// TODO: 上位三つの早いものを取り出す処理
	fastThree := domain.ApproachInfos{}
	// FIXME: 一旦それっぽいデータを埋めてる
	fastThree.ApproachInfo = append(fastThree.ApproachInfo, domain.ApproachInfo{})
	fastThree.ApproachInfo = append(fastThree.ApproachInfo, domain.ApproachInfo{})
	fastThree.ApproachInfo = append(fastThree.ApproachInfo, domain.ApproachInfo{})
	return fastThree
}



