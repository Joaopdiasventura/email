package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var registry = prometheus.NewRegistry()

var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "email",
		Subsystem: "api",
		Name:      "http_requests_total",
		Help:      "Total HTTP requests handled by the Vercel serverless function.",
	},
	[]string{"method", "route", "status"},
)

var httpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "email",
		Subsystem: "api",
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request duration for the Vercel serverless function.",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	},
	[]string{"method", "route", "status"},
)

var emailSendsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "email",
		Subsystem: "smtp",
		Name:      "sends_total",
		Help:      "Total email send attempts by result.",
	},
	[]string{"result"},
)

var smtpConfigErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "email",
		Subsystem: "smtp",
		Name:      "config_errors_total",
		Help:      "Total SMTP configuration load errors.",
	},
	[]string{"error"},
)

var metricsHandler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

func init() {
	registry.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		emailSendsTotal,
		smtpConfigErrorsTotal,
	)
}

func MetricsHandler() http.Handler {
	return metricsHandler
}

func RecordHTTPRequest(method string, route string, status int, duration time.Duration) {
	statusLabel := strconv.Itoa(status)

	httpRequestsTotal.WithLabelValues(method, route, statusLabel).Inc()
	httpRequestDuration.WithLabelValues(method, route, statusLabel).Observe(duration.Seconds())
}

func RecordEmailSend(result string) {
	emailSendsTotal.WithLabelValues(result).Inc()
}

func RecordSMTPConfigError(err string) {
	smtpConfigErrorsTotal.WithLabelValues(err).Inc()
}
