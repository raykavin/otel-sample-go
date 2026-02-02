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
	status       int
	forcedStatus int
	wroteHeader  bool
}

func (w *statusWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}

	if w.forcedStatus != 0 {
		code = w.forcedStatus
	}

	w.status = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		if w.forcedStatus != 0 {
			w.WriteHeader(w.forcedStatus)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
	return w.ResponseWriter.Write(b)
}

func Instrument(route string, h http.Handler, m *metrics.AppMetrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w}
		start := time.Now()

		if code, ok := statusFromRequest(r); ok {
			sw.forcedStatus = code
		}

		h.ServeHTTP(sw, r)

		if sw.status == 0 {
			if sw.forcedStatus != 0 {
				sw.status = sw.forcedStatus
			} else {
				sw.status = http.StatusOK
			}
		}

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
