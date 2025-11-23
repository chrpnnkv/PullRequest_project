package api

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var body ErrorBody
	body.Error.Code = code
	body.Error.Message = message

	_ = json.NewEncoder(w).Encode(body)
}
