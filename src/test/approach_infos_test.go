package test

import (
	"fmt"
	"github.com/shun-shun123/bus-timer/src/domain"
	"testing"
)

func TestGetFastThreePrune(t *testing.T) {
	approachInfos := domain.CreateApproachInfos()
	insertApproachInfo(&approachInfos, createApproachInfo(12, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(13, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(14, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(15, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(16, 0))

	fastThree := approachInfos.GetFastThree()

	checkLenOfSlice(fastThree, 3, t)
	checkRealArrivalTime(0, fastThree, "12:00", t)
	checkRealArrivalTime(1, fastThree, "13:00", t)
	checkRealArrivalTime(2, fastThree, "14:00", t)
}

func TestGetFastThreeSwitchOrder(t *testing.T) {
	approachInfos := domain.CreateApproachInfos()
	insertApproachInfo(&approachInfos, createApproachInfo(15, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(14, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(13, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(12, 0))
	insertApproachInfo(&approachInfos, createApproachInfo(11, 0))

	fastThree := approachInfos.GetFastThree()

	checkLenOfSlice(fastThree, 3, t)
	checkRealArrivalTime(0, fastThree, "11:00", t)
	checkRealArrivalTime(1, fastThree, "12:00", t)
	checkRealArrivalTime(2, fastThree, "13:00", t)

}

// 一つの接近情報を挿入する
func insertApproachInfo(approachInfos *domain.ApproachInfos, info domain.ApproachInfo) {
	approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, info)
}

func createApproachInfo(hour, min int) domain.ApproachInfo {
	info := domain.ApproachInfo{
		RealArrivalTime: fmt.Sprintf("%02d:%02d", hour, min),
	}
	return info
}

func checkRealArrivalTime(index int, fastThree domain.ApproachInfos, ans string, t *testing.T) {
	if fastThree.ApproachInfo[index].RealArrivalTime != ans {
		t.Fatalf("%v番目が%vになっています。", index, fastThree.ApproachInfo[index].RealArrivalTime)
	}
}

func checkLenOfSlice(fastThree domain.ApproachInfos, ans int, t *testing.T) {
	if len(fastThree.ApproachInfo) != ans {
		t.Fatalf("取得数が%vになっています。", len(fastThree.ApproachInfo))
	}
}