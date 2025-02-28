package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/timurzdev/mentorship-test-task/internal/deps"
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

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func New(metrics deps.Metrics) *Middleware {
	_ = metrics.RegisterHistogram(httpDurationMetricKey, httpDurationMetric)
	return &Middleware{metrics}
}

func (mw *Middleware) Handle(next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerDuration(httpDurationMetric, next)
}
