package test

import (
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"testing"
)

func TestScrapeApproachInfo(t *testing.T) {
	fetcher := infrastructure.ApproachInfoFetcher{}
	approachInfoUrl := "http://localhost:8000/rits_to_minakusa_enough_bus.html"
	viaUrl := ""
	approachInfos := fetcher.FetchApproachInfos(approachInfoUrl, viaUrl)
	if len(approachInfos.ApproachInfo) == 0 {
		t.Fatal("正しくスクレピングできていません.")
	}
	if len(approachInfos.ApproachInfo) != 3 {
		t.Fatalf("ApproachInfoの数が正しくありません。%v", len(approachInfos.ApproachInfo))
	}

	approachInfo := approachInfos.ApproachInfo[0]
	if approachInfo.MoreMin != "約8分後に到着" {
		t.Fatalf("[0]MoreMinが正しくありません。%v", approachInfo.MoreMin)
	}

	if approachInfo.RealArrivalTime != "14:25"{
		t.Fatalf("[0]RealArrivalTimeが正しくありません。%v", approachInfo.RealArrivalTime)
	}

	if approachInfo.Direction != "南草津駅" {
		t.Fatalf("[0]Directionが正しくありません。%v", approachInfo.Direction)
	}

	if approachInfo.ScheduledTime != "14:25" {
		t.Fatalf("[0]ScheduledTimeが正しくありません。%v", approachInfo.ScheduledTime)
	}

	if approachInfo.Delay != "(定時運行)" {
		t.Fatalf("[0].Delayが正しくありません。%v", approachInfo.Delay)
	}

	approachInfo = approachInfos.ApproachInfo[1]
	if approachInfo.MoreMin != "約18分後に到着" {
		t.Fatalf("[0]MoreMinが正しくありません。%v", approachInfo.MoreMin)
	}

	if approachInfo.RealArrivalTime != "14:35"{
		t.Fatalf("[0]RealArrivalTimeが正しくありません。%v", approachInfo.RealArrivalTime)
	}

	if approachInfo.Direction != "南草津駅" {
		t.Fatalf("[0]Directionが正しくありません。%v", approachInfo.Direction)
	}

	if approachInfo.ScheduledTime != "14:35" {
		t.Fatalf("[0]ScheduledTimeが正しくありません。%v", approachInfo.ScheduledTime)
	}

	if approachInfo.Delay != "(定時運行)" {
		t.Fatalf("[0].Delayが正しくありません。%v", approachInfo.Delay)
	}

	approachInfo = approachInfos.ApproachInfo[2]
	if approachInfo.MoreMin != "約37分後に到着" {
		t.Fatalf("[0]MoreMinが正しくありません。%v", approachInfo.MoreMin)
	}

	if approachInfo.RealArrivalTime != "14:54"{
		t.Fatalf("[0]RealArrivalTimeが正しくありません。%v", approachInfo.RealArrivalTime)
	}

	if approachInfo.Direction != "南草津駅" {
		t.Fatalf("[0]Directionが正しくありません。%v", approachInfo.Direction)
	}

	if approachInfo.ScheduledTime != "14:54" {
		t.Fatalf("[0]ScheduledTimeが正しくありません。%v", approachInfo.ScheduledTime)
	}

	if approachInfo.Delay != "(定時運行)" {
		t.Fatalf("[0].Delayが正しくありません。%v", approachInfo.Delay)
	}
}

func TestScrapeApproachInfoLessBus(t *testing.T) {
	fetcher := infrastructure.ApproachInfoFetcher{}
	approachInfoUrl := "http://localhost:8000/minakusa_to_rits_less_bus.html"
	viaUrl := ""
	approachInfos := fetcher.FetchApproachInfos(approachInfoUrl, viaUrl)
	if len(approachInfos.ApproachInfo) == 0 {
		t.Fatal("正しくスクレピングできていません.")
	}
	if len(approachInfos.ApproachInfo) != 1 {
		t.Fatalf("ApproachInfoの数が正しくありません。%v", len(approachInfos.ApproachInfo))
	}

	approachInfo := approachInfos.ApproachInfo[0]
	if approachInfo.MoreMin != "約28分後に到着" {
		t.Fatalf("[0]MoreMinが正しくありません。%v", approachInfo.MoreMin)
	}

	if approachInfo.RealArrivalTime != "17:25"{
		t.Fatalf("[0]RealArrivalTimeが正しくありません。%v", approachInfo.RealArrivalTime)
	}

	if approachInfo.Direction != "立命館大学" {
		t.Fatalf("[0]Directionが正しくありません。%v", approachInfo.Direction)
	}

	if approachInfo.ScheduledTime != "17:25" {
		t.Fatalf("[0]ScheduledTimeが正しくありません。%v", approachInfo.ScheduledTime)
	}

	if approachInfo.Delay != "(定時運行)" {
		t.Fatalf("[0].Delayが正しくありません。%v", approachInfo.Delay)
	}
}

func TestScrapeApproachInfoNoBus(t *testing.T) {
	fetcher := infrastructure.ApproachInfoFetcher{}
	approachInfoUrl := "http://localhost:8000/None.html"
	viaUrl := ""
	approachInfos := fetcher.FetchApproachInfos(approachInfoUrl, viaUrl)
	if len(approachInfos.ApproachInfo) != 0 {
		t.Fatalf("ApproachInfoの数が正しくありません。%v", len(approachInfos.ApproachInfo))
	}

}