package metrics

import (
	"context"
	"runtime"
	"sync"
	"telemetry-sample/internal/system"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func RegisterRuntimeMetrics(m *AppMetrics) {
	var (
		lastCPUSeconds float64
		lastWall       time.Time
		mu             sync.Mutex
	)

	cbFn := func(ctx context.Context, o metric.Observer) error {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		o.ObserveFloat64(m.MemGauge, float64(mem.Alloc))
		o.ObserveInt64(m.GoroutinesGauge, int64(runtime.NumGoroutine()))

		cpuSeconds, err := system.GetCPUSeconds()
		if err != nil {
			return nil
		}

		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		if !lastWall.IsZero() {
			deltaCPU := cpuSeconds - lastCPUSeconds
			deltaWall := now.Sub(lastWall).Seconds()

			if deltaWall > 0 {
				cpuPercent := (deltaCPU / deltaWall) * 100 / float64(runtime.NumCPU())
				if cpuPercent > 0 {
					o.ObserveFloat64(m.CPUGauge, cpuPercent)
				}
			}
		}

		lastCPUSeconds = cpuSeconds
		lastWall = now
		return nil
	}

	if _, err := otel.Meter("go-app").RegisterCallback(cbFn,
		m.MemGauge, m.GoroutinesGauge, m.CPUGauge); err != nil {
		panic(err)
	}
}
