package test

import (
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"testing"
)

func TestSample(t *testing.T) {
	var err error = nil
	if err != nil {
		t.Fatalf("failed test %v", err)
	}
}

func TestScrapeApproachInfo(t *testing.T) {
	fetcher := infrastructure.ApproachInfoFetcher{}
	approachInfoUrl := "http://localhost:8000/rits_to_minakusa_enough_bus.html"
	viaUrl := ""
	approachInfos := fetcher.FetchApproachInfos(approachInfoUrl, viaUrl)
	if len(approachInfos.ApproachInfo) == 0 {
		t.Fatal("正しくスクレピングできていません.")
	}
	for i, v := range approachInfos.ApproachInfo {
		t.Log(i, ".moreMin: ", v.MoreMin)
		t.Log(i, ".RealArrivalTime: ", v.RealArrivalTime)
		t.Log(i, ".Direction: ", v.Direction)
		t.Log(i, ".Via: ", v.Via)
		t.Log(i, ".ScheduledTime: ", v.ScheduledTime)
		t.Log(i, ".Delay: ", v.Delay)
		t.Log(i, ".BusStop: ", v.BusStop)
	}
}