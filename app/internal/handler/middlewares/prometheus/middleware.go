package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/PKSonlem/micro/internal/deps"
)

var (
	httpDurationMetricKey = "httpDurationMetric"
	httpDurationMetric    = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_ms",
			Help:    "Duration of HTTP requests in millisecons.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"code", "method"},
	)
)

type Middleware struct {
	metrics deps.Metrics
}

func New(metrics deps.Metrics) *Middleware {
	_ = metrics.RegisterHistogram(httpDurationMetricKey, httpDurationMetric)
	return &Middleware{metrics}
}

func (mw *Middleware) Handle(next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerDuration(httpDurationMetric, next)
}
