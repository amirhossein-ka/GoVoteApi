package mux

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
