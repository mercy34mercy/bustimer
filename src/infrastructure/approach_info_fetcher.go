package infrastructure

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/shun-shun123/bus-timer/src/domain"
)

type ApproachInfoFetcher struct {
}

func (fetcher ApproachInfoFetcher) FetchApproachInfos(url string) domain.ApproachInfos {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		// TODO: スクレイピングが失敗した場合の処理
		// 近江鉄道バスサーバ死亡説...
		return domain.ApproachInfos{}
	}
	fmt.Println(doc)
	approachInfos := domain.ApproachInfos{}
	// さぁ、スクレイピングと行こうか。
	// TODO: スクレイピングが成功したら情報を取り出す
	return approachInfos
}
