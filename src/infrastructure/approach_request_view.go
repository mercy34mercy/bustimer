package infrastructure

import (
	"net/http"

	"fmt"
	"time"

	"github.com/shun-shun123/bus-timer/src/presenter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("bustimer/approach")

func ApproachInfoRequest(c Context) error {
	// トレーシング開始
	ctx, span := tracer.Start(c.Request(), "ApproachInfoRequest")
	defer span.End()

	startTime := time.Now()

	// リクエストURLを作成
	approachInfoUrls := c.GetApproachInfoUrls()
	span.SetAttributes(
		attribute.StringSlice("approachInfoUrls", approachInfoUrls),
		attribute.Int("urlCount", len(approachInfoUrls)),
		attribute.String("endpoint", "/bus/time/v3"),
		attribute.String("method", "GET"),
	)

	from, to := c.GetFromToQuery()
	span.SetAttributes(
		attribute.Int("from", int(from)),
		attribute.Int("to", int(to)),
		attribute.String("fromString", from.String()),
		attribute.String("toString", to.String()),
	)

	// URLからデータを取得するFetcherを作成
	fetcher := ApproachInfoFetcher{
		from: from,
		to:   to,
	}

	// サブスパンを作成してデータ取得部分を計測
	_, fetchSpan := tracer.Start(ctx, "fetchApproachData")
	approachInfos := presenter.RequestApproachInfosWithContext(ctx, approachInfoUrls, fetcher)
	fetchSpan.SetAttributes(
		attribute.Int("resultCount", len(approachInfos.ApproachInfo)),
		attribute.String("urlsFetched", fmt.Sprintf("%v", approachInfoUrls)),
	)
	fetchSpan.End()

	// 処理時間の計測
	processingTime := time.Since(startTime)
	span.SetAttributes(
		attribute.String("processingTime", processingTime.String()),
		attribute.Int64("processingTimeMs", processingTime.Milliseconds()),
	)

	// レスポンス処理
	if len(approachInfos.ApproachInfo) == 0 {
		span.SetAttributes(
			attribute.String("status", "noContent"),
			attribute.Int("httpStatus", http.StatusNoContent),
		)
		span.SetStatus(codes.Error, "No approach info found")
		return c.Response("ApproachInfoRequest", http.StatusNoContent, approachInfos)
	}

	span.SetAttributes(
		attribute.String("status", "success"),
		attribute.Int("httpStatus", http.StatusOK),
		attribute.Int("resultCount", len(approachInfos.ApproachInfo)),
	)
	return c.Response("ApproachInfoRequest", http.StatusOK, approachInfos)
}

func ApproachInfoRequestV2(c Context) error {
	// トレーシング開始
	ctx, span := tracer.Start(c.Request(), "ApproachInfoRequestV2")
	defer span.End()

	startTime := time.Now()

	// リクエストURLを作成
	approachInfoUrls := c.GetApproachInfoUrls()
	span.SetAttributes(
		attribute.StringSlice("approachInfoUrls", approachInfoUrls),
		attribute.Int("urlCount", len(approachInfoUrls)),
		attribute.String("endpoint", "/bus/time/v2"),
		attribute.String("method", "GET"),
	)

	from, to := c.GetFromToQuery()
	span.SetAttributes(
		attribute.Int("from", int(from)),
		attribute.Int("to", int(to)),
		attribute.String("fromString", from.String()),
		attribute.String("toString", to.String()),
	)

	// URLからデータを取得するFetcherを作成
	fetcher := ApproachInfoFetcher{
		from: from,
		to:   to,
	}

	// サブスパンを作成してデータ取得部分を計測
	_, fetchSpan := tracer.Start(ctx, "fetchApproachDataV2")
	approachInfos := presenter.RequestApproachInfosWithContext(ctx, approachInfoUrls, fetcher)
	fetchSpan.SetAttributes(
		attribute.Int("resultCount", len(approachInfos.ApproachInfo)),
		attribute.String("urlsFetched", fmt.Sprintf("%v", approachInfoUrls)),
	)
	fetchSpan.End()

	// V2形式に変換
	v2 := approachInfos.CopyToV2()
	span.SetAttributes(attribute.Int("v2ResultCount", len(v2.ApproachInfo)))

	// 処理時間の計測
	processingTime := time.Since(startTime)
	span.SetAttributes(
		attribute.String("processingTime", processingTime.String()),
		attribute.Int64("processingTimeMs", processingTime.Milliseconds()),
	)

	// レスポンス処理
	if len(approachInfos.ApproachInfo) == 0 {
		span.SetAttributes(
			attribute.String("status", "noContent"),
			attribute.Int("httpStatus", http.StatusNoContent),
		)
		span.SetStatus(codes.Error, "No approach info found")
		return c.Response("ApproachInfoRequestV2", http.StatusNoContent, v2)
	}

	span.SetAttributes(
		attribute.String("status", "success"),
		attribute.Int("httpStatus", http.StatusOK),
		attribute.Int("resultCount", len(v2.ApproachInfo)),
	)
	return c.Response("ApproachInfoRequestV2", http.StatusOK, v2)
}
