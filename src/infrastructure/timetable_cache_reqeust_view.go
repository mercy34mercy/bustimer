package infrastructure

import (
	"fmt"
	"github.com/shun-shun123/bus-timer/src/domain"
	"time"
)

const (
	KeyRits = "Rits"
	KeyMinakusa  = "Minakusa"
)

var TimeTable = make(map[string]domain.TimeTable)

func TimeTableCacheStart() {
	for ;; {
		fmt.Println("時刻表の更新をします")
		fetcher := TimetableFetcher{}
		TimeTable[KeyRits] = fetcher.FetchTimetable(FrRits)
		TimeTable[KeyMinakusa] = fetcher.FetchTimetable(FrMinakusa)
		time.Sleep(30 * time.Second)
	}
}