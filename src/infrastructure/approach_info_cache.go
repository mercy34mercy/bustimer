package infrastructure

import (
	"sync"
	"time"

	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
)

// キャッシュのエントリー
type ApproachInfoCacheEntry struct {
	Data      domain.ApproachInfos
	CreatedAt time.Time
}

// frとtoの組み合わせをキーとするキャッシュマップ
var (
	approachInfoCache    = make(map[string]ApproachInfoCacheEntry)
	approachInfoCacheMtx sync.RWMutex
	approachInfoCacheTTL = 1 * time.Minute
)

// キャッシュのキーを生成する
func createCacheKey(from config.From, to config.To) string {
	return from.String() + "-" + to.String()
}

// キャッシュからデータを取得する
func GetApproachInfoFromCache(from config.From, to config.To) (domain.ApproachInfos, bool) {
	approachInfoCacheMtx.RLock()
	defer approachInfoCacheMtx.RUnlock()

	key := createCacheKey(from, to)
	entry, exists := approachInfoCache[key]
	if !exists {
		return domain.ApproachInfos{}, false
	}

	// TTLをチェック
	if time.Since(entry.CreatedAt) > approachInfoCacheTTL {
		// 有効期限切れの場合は存在しないものとして扱う
		return domain.ApproachInfos{}, false
	}

	return entry.Data, true
}

// キャッシュにデータを保存する
func SetApproachInfoToCache(from config.From, to config.To, data domain.ApproachInfos) {
	approachInfoCacheMtx.Lock()
	defer approachInfoCacheMtx.Unlock()

	key := createCacheKey(from, to)
	approachInfoCache[key] = ApproachInfoCacheEntry{
		Data:      data,
		CreatedAt: time.Now(),
	}
}
