package infrastructure

import (
	"context"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

// InitTracer はCloud Traceの初期化を行います
func InitTracer() (*sdktrace.TracerProvider, error) {
	// Cloud Runでのデフォルト認証情報を使用するとGCPへの認証が自動的に行われます
	ctx := context.Background()

	// サービス名を取得（環境変数から）
	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		serviceName = "bustimer-service"
	}

	// リソース情報を設定
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Cloud Traceエクスポーターを作成
	exporter, err := createExporter(ctx)
	if err != nil {
		return nil, err
	}

	// トレースプロバイダーを設定
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)

	// コンテキスト伝播の設定
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

// createExporter はCloud Traceエクスポーターを作成します
func createExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	// Cloud Run環境ではデフォルト認証情報を使用
	opts := []option.ClientOption{}
	creds, err := oauth.NewApplicationDefault(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Printf("Failed to create credentials: %v", err)
		// 認証に失敗してもエクスポーターの作成は続行
	} else {
		opts = append(opts, option.WithGRPCDialOption(grpc.WithPerRPCCredentials(creds)))
	}

	// Cloud Traceエンドポイント
	endpoint := os.Getenv("TRACE_EXPORTER_ENDPOINT")
	if endpoint == "" {
		endpoint = "cloudtrace-exporter.googleapis.com:443"
	}

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithTimeout(5*time.Second),
	)

	return otlptrace.New(ctx, client)
}

// ShutdownTracer はトレースプロバイダーをシャットダウンします
func ShutdownTracer(ctx context.Context, tp *sdktrace.TracerProvider) {
	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}
}
