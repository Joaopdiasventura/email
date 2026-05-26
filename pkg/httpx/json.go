package httpx

import (
	"encoding/json"
	"net/http"
)

type Response map[string]any

func WriteJSON(w http.ResponseWriter, status int, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
