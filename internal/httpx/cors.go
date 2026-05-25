package httpx

import "net/http"

var allowedOrigins = map[string]bool{
	"http://localhost:4200":            true,
	"https://joaopdias-dev.vercel.app": true,
	"https://joaopdias.dev.br":         true,
}

func ApplyCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
