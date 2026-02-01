package server

import (
	"context"
	"log"
	"math/rand"
	"net/http"
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
