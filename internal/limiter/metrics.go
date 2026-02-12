package limiter

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	limiterMax = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "adaptive_max_inflight",
			Help: "Current max inflight requests",
		},
	)

	limiterEMA = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "adaptive_latency_ema_ms",
			Help: "EMA latency in ms",
		},
	)
)

func Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		limiterMax,
		limiterEMA,
	}
}

func (l *AdaptiveLimiter) ExportMetrics() {
	limiterMax.Set(float64(l.maxInFlight.Load()))
	limiterEMA.Set(float64(time.Duration(l.emaLatency.Load()).Milliseconds()))
}
