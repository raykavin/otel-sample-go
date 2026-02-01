package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type AppMetrics struct {
	RequestCount    metric.Int64Counter
	RequestDuration metric.Float64Histogram
	MemGauge        metric.Float64ObservableGauge
	GoroutinesGauge metric.Int64ObservableGauge
	CPUGauge        metric.Float64ObservableGauge
}

func Setup() *AppMetrics {
	meter := otel.Meter("go-app")

	mem, _ := meter.Float64ObservableGauge("process_memory_heap_bytes")
	goroutines, _ := meter.Int64ObservableGauge("process_goroutines")
	cpu, _ := meter.Float64ObservableGauge("process_cpu_percent")

	reqCount, _ := meter.Int64Counter("http_server_requests_total")
	reqDuration, _ := meter.Float64Histogram("http_server_request_duration_ms")

	return &AppMetrics{
		RequestCount:    reqCount,
		RequestDuration: reqDuration,
		MemGauge:        mem,
		GoroutinesGauge: goroutines,
		CPUGauge:        cpu,
	}
}
