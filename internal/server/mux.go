package server

import (
	"net/http"
	"telemetry-sample/internal/metrics"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func BuildMux(m *metrics.AppMetrics) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(indexHandler))

	mux.Handle("/hello", Instrument(
		"/hello",
		otelhttp.NewHandler(helloHandler(), "HelloEndpoint"),
		m,
	))

	mux.Handle("/cpu", Instrument(
		"/cpu",
		otelhttp.NewHandler(cpuHandler(), "CpuEndpoint"),
		m,
	))

	mux.Handle("/memory", Instrument(
		"/memory",
		otelhttp.NewHandler(memoryHandler(), "MemoryEndpoint"),
		m,
	))

	return mux
}
