package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
	"time"
)

var TimeTable = make(map[config.From]domain.TimeTable)

// 時刻表データのキャッシュを作成する
// goroutineで実行すれば、configで設定されているタイマーの頻度でデータを更新する
func TimeTableCacheStart() {
	fetcher := TimetableFetcher{}
	for {
		TimeTable[config.FromRits] = fetcher.FetchTimetable(config.FromRits, config.ToMinakusa)
		TimeTable[config.FromMinakusa] = fetcher.FetchTimetable(config.FromMinakusa, config.ToRits)
		time.Sleep(config.TimeTableCacheUpdateDuration)
	}
}
