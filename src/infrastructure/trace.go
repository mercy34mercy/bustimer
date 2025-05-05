package infrastructure

import (
	"context"
	"log"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

// InitTracer はCloud Traceの初期化を行います
// 呼び出し側はtp.Shutdown()でクリーンアップを行う必要があります
func InitTracer() (*sdktrace.TracerProvider, error) {
	log.Println("Starting Cloud Trace initialization...")
	ctx := context.Background()

	// サービス名を取得（環境変数から）
	serviceName := os.Getenv("K_SERVICE")
	if serviceName == "" {
		serviceName = "bustimer-service"
		log.Println("K_SERVICE environment variable not set, using default service name:", serviceName)
	} else {
		log.Println("Using service name from K_SERVICE:", serviceName)
	}

	// プロジェクトIDを環境変数から取得
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		projectID = "your-project-id" // デフォルト値は適切なプロジェクトIDに変更してください
		log.Println("PROJECT_ID environment variable not set, using default project ID:", projectID)
	} else {
		log.Println("Using project ID from environment variable:", projectID)
	}

	// Cloud Traceエクスポーターを作成
	log.Println("Creating Cloud Trace exporter...")
	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		log.Printf("Failed to create exporter: %v", err)
		return nil, err
	}
	log.Println("Cloud Trace exporter created successfully")

	// リソース情報を設定
	log.Println("Creating resource with service name:", serviceName)
	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		log.Printf("Failed to create resource: %v", err)
		return nil, err
	}
	log.Println("Resource created successfully")

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

// ShutdownTracer はトレースプロバイダーをシャットダウンします
func ShutdownTracer(ctx context.Context, tp *sdktrace.TracerProvider) {
	log.Println("Shutting down tracer provider...")
	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	} else {
		log.Println("Tracer provider shut down successfully")
	}
}
