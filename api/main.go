package api

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Joaopdiasventura/email/pkg/config"
	"github.com/Joaopdiasventura/email/pkg/httpx"
	"github.com/Joaopdiasventura/email/pkg/mail"
	"github.com/Joaopdiasventura/email/pkg/observability"
)

const (
	contactEmail = "joaopdias.dev@gmail.com"
	contactRoute = "/contact"
	metricsRoute = "/metrics"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	recorder := &statusRecorder{ResponseWriter: w}
	start := time.Now()
	route := routeLabel(r)
	result := "ok"

	logger := observability.Logger().With(
		"request_id", requestID(r),
		"method", r.Method,
		"path", r.URL.Path,
		"route", route,
	)

	defer func() {
		duration := time.Since(start)
		status := recorder.Status()

		observability.RecordHTTPRequest(r.Method, route, status, duration)
		logger.Info("request_completed",
			"status", status,
			"result", result,
			"duration_ms", float64(duration.Microseconds())/1000,
		)
	}()

	httpx.ApplyCORS(recorder, r)

	if r.Method == http.MethodOptions {
		recorder.WriteHeader(http.StatusNoContent)
		return
	}

	if isMetricsRequest(r) {
		if r.Method != http.MethodGet {
			result = "method_not_allowed"
			recorder.Header().Set("Allow", "GET")
			httpx.WriteJSON(recorder, http.StatusMethodNotAllowed, httpx.Response{
				"error": "METHOD_NOT_ALLOWED",
			})
			return
		}

		if !metricsAuthorized(r) {
			result = "metrics_unauthorized"
			logger.Warn("metrics_rejected", "reason", result)
			httpx.WriteJSON(recorder, http.StatusUnauthorized, httpx.Response{
				"error": "UNAUTHORIZED",
			})
			return
		}

		observability.MetricsHandler().ServeHTTP(recorder, r)
		return
	}

	if r.Method != http.MethodPost {
		result = "method_not_allowed"
		logger.Warn("request_rejected", "reason", result)
		recorder.Header().Set("Allow", "POST")
		httpx.WriteJSON(recorder, http.StatusMethodNotAllowed, httpx.Response{
			"error": "METHOD_NOT_ALLOWED",
		})
		return
	}

	var body mail.Request

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		result = "invalid_json"
		logger.Warn("request_rejected", "reason", result, "error", err.Error())
		httpx.WriteJSON(recorder, http.StatusBadRequest, httpx.Response{
			"error": "INVALID_JSON",
		})
		return
	}

	body.Normalize()

	if !body.Valid() {
		result = "invalid_payload"
		logger.Warn("request_rejected",
			"reason", result,
			"has_text", body.Text != "",
			"has_html", body.HTML != "",
		)
		httpx.WriteJSON(recorder, http.StatusBadRequest, httpx.Response{
			"error": "INVALID_PAYLOAD",
		})
		return
	}

	smtpConfig, err := config.LoadSMTPConfig()
	if err != nil {
		result = "smtp_config_error"
		observability.RecordSMTPConfigError(err.Error())
		logger.Error("smtp_config_load_failed", "error", err.Error())
		httpx.WriteJSON(recorder, http.StatusInternalServerError, httpx.Response{
			"error": err.Error(),
		})
		return
	}

	sender := mail.NewSender(smtpConfig)

	if err := sender.Send(contactEmail, body); err != nil {
		result = "send_failed"
		observability.RecordEmailSend("failure")
		logger.Error("email_send_failed", "error", err.Error())
		httpx.WriteJSON(recorder, http.StatusInternalServerError, httpx.Response{
			"error": "SEND_FAILED",
		})
		return
	}

	observability.RecordEmailSend("success")
	logger.Info("email_sent",
		"has_text", body.Text != "",
		"has_html", body.HTML != "",
		"reply_to_provided", body.ReplyTo != "",
	)

	httpx.WriteJSON(recorder, http.StatusOK, httpx.Response{
		"ok": true,
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}

	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}

	return r.ResponseWriter.Write(body)
}

func (r *statusRecorder) Status() int {
	if r.status == 0 {
		return http.StatusOK
	}

	return r.status
}

func isMetricsRequest(r *http.Request) bool {
	return r.URL.Path == metricsRoute || r.URL.Query().Get("metrics") == "1"
}

func metricsAuthorized(r *http.Request) bool {
	token := os.Getenv("METRICS_TOKEN")
	if token == "" {
		return true
	}

	expected := "Bearer " + token
	actual := r.Header.Get("Authorization")

	return subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) == 1
}

func routeLabel(r *http.Request) string {
	if isMetricsRequest(r) {
		return metricsRoute
	}

	return contactRoute
}

func requestID(r *http.Request) string {
	for _, header := range []string{"X-Vercel-Id", "X-Request-Id", "X-Correlation-Id"} {
		value := r.Header.Get(header)
		if value != "" {
			return value
		}
	}

	return ""
}
