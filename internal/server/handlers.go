package server

import (
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Endpoints: /hello /cpu /memory /metrics"))
}

func helloHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanLog(r.Context(), "hello request received")

		time.Sleep(getRandomDur())

		_, _ = w.Write([]byte("Hello OpenTelemetry"))
	})
}

func cpuHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanLog(r.Context(), "cpu work requested")

		start := time.Now()
		for time.Since(start) < getRandomDur() {
			_ = time.Now().UnixNano()
		}

		_, _ = w.Write([]byte("CPU work done"))
	})
}

func memoryHandler() http.Handler {
	var blobs [][]byte

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanLog(r.Context(), "memory allocation requested")

		blobs = append(blobs, make([]byte, 5*1024*1024))
		if len(blobs) > 5 {
			blobs = blobs[1:]
		}

		time.Sleep(getRandomDur())

		_, _ = w.Write([]byte("Allocated ~5MB"))
	})
}
