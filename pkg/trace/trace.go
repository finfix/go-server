package trace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"pkg/errors"
)

type TracerConfig struct {
	OtlHTTPEndpoint string `env:"OTLP_HTTP_ENDPOINT"`
	OtlInsecure     bool   `env:"OTLP_INSECURE"`
}

func StartTracing(cfg TracerConfig, serviceName string) error {
	headers := map[string]string{
		"content-type": "application/json",
	}

	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.OtlHTTPEndpoint),
		otlptracehttp.WithHeaders(headers),
	}

	if cfg.OtlInsecure {
		options = append(options, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(options...),
	)
	if err != nil {
		return errors.InternalServer.Wrap(err)
	}

	tracerprovider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)

	otel.SetTracerProvider(tracerprovider)

	logExporter, err := otlploghttp.New(context.Background(),
		otlploghttp.WithEndpoint(cfg.OtlHTTPEndpoint),
		otlploghttp.WithHeaders(headers),
		otlploghttp.WithInsecure(),
	)

	// Create the logger provider
	lp := log.NewLoggerProvider(
		log.WithProcessor(
			log.NewBatchProcessor(logExporter),
		),
	)

	global.SetLoggerProvider(lp)

	return nil
}
