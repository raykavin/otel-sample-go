package telemetry

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const shutdownTimeout = 5 * time.Second

type config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	OTLPInsecure   bool
}

// loadConfig reads telemetry configuration from environment variables
func loadConfig() config {
	return config{
		ServiceName:    getenv("OTEL_SERVICE_NAME", "go-telemetry-app"),
		ServiceVersion: getenv("SERVICE_VERSION", "dev"),
		Environment:    getenv("APP_ENV", "local"),
		OTLPEndpoint:   getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		OTLPInsecure:   strings.ToLower(getenv("OTEL_EXPORTER_OTLP_INSECURE", "true")) != "false",
	}
}

// newResource creates an OpenTelemetry resource with service metadata
func newResource(ctx context.Context, cfg config) (*resource.Resource, error) {
	return resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", cfg.ServiceVersion),
			attribute.String("deployment.environment", cfg.Environment),
		),
	)
}

// initTracing configures the OTLP trace exporter and tracer provider
func initTracing(
	ctx context.Context,
	res *resource.Resource,
	cfg config,
) (*sdktrace.TracerProvider, error) {
	options := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
	}

	if cfg.OTLPInsecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}

// initMetrics configures the Prometheus exporter and HTTP endpoint
func initMetrics(res *resource.Resource) error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exporter),
	)

	otel.SetMeterProvider(meterProvider)

	// Expose Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	return nil
}

// getenv returns the environment variable value or a default if empty
func getenv(key, def string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return def
	}
	return value
}

// InitTelemetry initializes tracing (OTLP) and metrics (Prometheus)
// and returns a shutdown function.
func InitTelemetry(ctx context.Context) (func(), error) {
	cfg := loadConfig()

	res, err := newResource(ctx, cfg)
	if err != nil {
		return nil, err
	}

	tracerProvider, err := initTracing(ctx, res, cfg)
	if err != nil {
		return nil, err
	}

	if err := initMetrics(res); err != nil {
		return nil, err
	}

	// Use W3C Trace Context propagation
	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Println("Telemetry initialized")

	return func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		_ = tracerProvider.Shutdown(shutdownCtx)
	}, nil
}
