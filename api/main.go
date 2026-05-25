package api

import (
	"encoding/json"
	"net/http"

	"github.com/Joaopdiasventura/email/internal/config"
	"github.com/Joaopdiasventura/email/internal/httpx"
	"github.com/Joaopdiasventura/email/internal/mail"
)

const contactEmail = "joaopdias.dev@gmail.com"

func Handler(w http.ResponseWriter, r *http.Request) {
	httpx.ApplyCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		httpx.WriteJSON(w, http.StatusMethodNotAllowed, httpx.Response{
			"error": "METHOD_NOT_ALLOWED",
		})
		return
	}

	var body mail.Request

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, httpx.Response{
			"error": "INVALID_JSON",
		})
		return
	}

	body.Normalize()

	if !body.Valid() {
		httpx.WriteJSON(w, http.StatusBadRequest, httpx.Response{
			"error": "INVALID_PAYLOAD",
		})
		return
	}

	smtpConfig, err := config.LoadSMTPConfig()
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.Response{
			"error": err.Error(),
		})
		return
	}

	sender := mail.NewSender(smtpConfig)

	if err := sender.Send(contactEmail, body); err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.Response{
			"error": "SEND_FAILED",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.Response{
		"ok": true,
	})
}
