package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	InFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_in_flight_requests",
		Help: "Current number of in-flight HTTP requests"})

	ShedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_shed_requests_total",
		Help: "Total number of shed HTTP requests",
	})

	RequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency",
		Buckets: prometheus.DefBuckets,
	})
)

func Register(extra ...prometheus.Collector) {
	prometheus.MustRegister(InFlight, ShedTotal, RequestDuration)

	for _, c := range extra {
		prometheus.MustRegister(c)
	}
}
