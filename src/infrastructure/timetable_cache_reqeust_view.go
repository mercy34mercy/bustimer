package infrastructure

import (
	"fmt"
	"github.com/shun-shun123/bus-timer/src/domain"
	"time"
)

var TimeTable = make(map[string]domain.TimeTable)

func TimeTableCacheStart() {
	for ;; {
		fmt.Println("時刻表の更新をします")
		fetcher := TimetableFetcher{}
		TimeTable[FrRits] = fetcher.FetchTimetable(FrRits)
		TimeTable[FrMinakusa] = fetcher.FetchTimetable(FrMinakusa)
		time.Sleep(30 * time.Second)
	}
}