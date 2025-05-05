package presenter

import (
	"context"

	"github.com/shun-shun123/bus-timer/src/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("bustimer/presenter")

func RequestApproachInfos(approachInfoUrls []string, fetcher IFetchApproachInfos) domain.ApproachInfos {
	return RequestApproachInfosWithContext(context.Background(), approachInfoUrls, fetcher)
}

func RequestApproachInfosWithContext(ctx context.Context, approachInfoUrls []string, fetcher IFetchApproachInfos) domain.ApproachInfos {
	_, span := tracer.Start(ctx, "RequestApproachInfos")
	defer span.End()

	approachInfos := domain.CreateApproachInfos()
	span.SetAttributes(attribute.Int("urlCount", len(approachInfoUrls)))

	for _, url := range approachInfoUrls {
		fetchCtx, fetchSpan := tracer.Start(ctx, "FetchApproachInfos")
		fetchSpan.SetAttributes(attribute.String("url", url))

		fetchResult := fetcher.FetchApproachInfosWithContext(fetchCtx, url, approachInfos)
		fetchSpan.SetAttributes(attribute.Int("fetchResultCount", len(fetchResult.ApproachInfo)))
		fetchSpan.End()

		approachInfos.ApproachInfo = append(approachInfos.ApproachInfo, fetchResult.ApproachInfo...)
	}

	fastThree := approachInfos.GetFastThree()
	span.SetAttributes(attribute.Int("fastThreeCount", len(fastThree.ApproachInfo)))

	return fastThree
}
