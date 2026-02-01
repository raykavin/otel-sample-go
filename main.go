package main

import (
	"context"
	"log"
	"net/http"

	"telemetry-sample/internal/logging"
	"telemetry-sample/internal/metrics"
	"telemetry-sample/internal/server"
	"telemetry-sample/pkg/telemetry"
)

func main() {
	ctx := context.Background()

	logging.Setup()

	shutdown, err := telemetry.InitTelemetry(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown()

	m := metrics.Setup()
	metrics.RegisterRuntimeMetrics(m)

	appMux := server.BuildMux(m)

	startServers(appMux)
}

func startServers(appMux *http.ServeMux) {
	go func() {
		log.Println("App running on :8080")
		log.Fatal(http.ListenAndServe(":8080", appMux))
	}()

	log.Println("Metrics running on :9464")
	log.Fatal(http.ListenAndServe(":9464", http.DefaultServeMux))
}
