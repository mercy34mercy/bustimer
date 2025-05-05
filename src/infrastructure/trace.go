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
	log.Println("Starting Cloud Trace initialization...")
	// Cloud Runでのデフォルト認証情報を使用するとGCPへの認証が自動的に行われます
	ctx := context.Background()

	// サービス名を取得（環境変数から）
	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		serviceName = "bustimer-service"
		log.Println("K_SERVICE environment variable not set, using default service name:", serviceName)
	} else {
		log.Println("Using service name from K_SERVICE:", serviceName)
	}

	// リソース情報を設定
	log.Println("Creating resource with service name:", serviceName)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		log.Printf("Failed to create resource: %v", err)
		return nil, err
	}
	log.Println("Resource created successfully")

	// Cloud Traceエクスポーターを作成
	log.Println("Creating Cloud Trace exporter...")
	exporter, err := createExporter(ctx)
	if err != nil {
		log.Printf("Failed to create exporter: %v", err)
		return nil, err
	}
	log.Println("Cloud Trace exporter created successfully")

	// トレースプロバイダーを設定
	log.Println("Setting up tracer provider...")
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	log.Println("Tracer provider configured successfully")

	// コンテキスト伝播の設定
	otel.SetTextMapPropagator(propagation.TraceContext{})
	log.Println("Text map propagator set to TraceContext")

	log.Println("Cloud Trace initialization completed successfully")
	return tp, nil
}

// createExporter はCloud Traceエクスポーターを作成します
func createExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	log.Println("Creating Cloud Trace exporter...")
	// Cloud Run環境ではデフォルト認証情報を使用
	opts := []option.ClientOption{}

	log.Println("Attempting to create OAuth credentials...")
	creds, err := oauth.NewApplicationDefault(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Printf("Failed to create credentials: %v", err)
		// 認証に失敗してもエクスポーターの作成は続行
	} else {
		log.Println("OAuth credentials created successfully")
		opts = append(opts, option.WithGRPCDialOption(grpc.WithPerRPCCredentials(creds)))
	}

	// Cloud Traceエンドポイント
	endpoint := os.Getenv("TRACE_EXPORTER_ENDPOINT")
	if endpoint == "" {
		endpoint = "cloudtrace-exporter.googleapis.com:443"
		log.Println("TRACE_EXPORTER_ENDPOINT not set, using default endpoint:", endpoint)
	} else {
		log.Println("Using Cloud Trace endpoint from environment variable:", endpoint)
	}

	log.Println("Creating OTLP trace gRPC client with endpoint:", endpoint)
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithTimeout(5*time.Second),
	)

	log.Println("Creating OTLP trace exporter...")
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Printf("Failed to create OTLP trace exporter: %v", err)
		return nil, err
	}
	log.Println("OTLP trace exporter created successfully")
	return exporter, nil
}

// ShutdownTracer はトレースプロバイダーをシャットダウンします
func ShutdownTracer(ctx context.Context, tp *sdktrace.TracerProvider) {
	log.Println("Shutting down tracer provider...")
	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	} else {
		log.Println("Tracer provider shut down successfully")
	}
}
