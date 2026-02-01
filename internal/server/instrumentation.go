package server

import (
	"net/http"
	"telemetry-sample/internal/metrics"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Instrument(route string, h http.Handler, m *metrics.AppMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: getRandStatusCode()}
		start := time.Now()

		h.ServeHTTP(sw, r)

		attrs := []attribute.KeyValue{
			attribute.String("http.route", route),
			attribute.String("http.method", r.Method),
			attribute.Int("http.status_code", sw.status),
		}

		m.RequestCount.Add(r.Context(), 1, metric.WithAttributes(attrs...))
		m.RequestDuration.Record(r.Context(),
			float64(time.Since(start).Milliseconds()),
			metric.WithAttributes(attrs...),
		)
	})
}
