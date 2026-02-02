package server

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

var (
	statusCodes = []int{
		http.StatusOK,
		http.StatusInternalServerError,
		http.StatusCreated,
		http.StatusBadRequest,
		http.StatusTemporaryRedirect,
		http.StatusNotFound,
		http.StatusBadGateway,
		http.StatusUnauthorized,
		http.StatusFound,
		http.StatusConflict,
		http.StatusUnprocessableEntity,
		http.StatusConflict,
		http.StatusPartialContent,
		http.StatusFailedDependency,
		http.StatusGatewayTimeout,
	}

	src    = rand.NewSource(time.Now().Unix())
	rd     = rand.New(src)
	maxDur = 500
	defDur = time.Millisecond
)

func spanLog(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()

	log.Printf("%s trace_id=%s", msg, traceID)
}

func getRandomDur() time.Duration {
	return time.Duration(rd.Intn(maxDur) * int(defDur))
}

func getRandStatusCode() int {
	idx := rd.Intn(len(statusCodes))
	return statusCodes[idx]
}

func statusFromRequest(r *http.Request) (int, bool) {
	if code, ok := parseStatus(r.URL.Query().Get("status")); ok {
		return code, true
	}

	if code, ok := parseStatus(r.Header.Get("X-Debug-Status")); ok {
		return code, true
	}

	if isTruthy(os.Getenv("SIMULATE_STATUS_CODES")) {
		return getRandStatusCode(), true
	}

	return 0, false
}

func parseStatus(value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}

	code, err := strconv.Atoi(value)
	if err != nil || code < 100 || code > 599 {
		return 0, false
	}

	return code, true
}

func isTruthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}
