package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondJSON sends a JSON response with the given status code and data.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(status)
	

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// This is a server-side error, so we log it and don't try to
			// send a new response body since headers have already been sent.
			log.Printf("failed to encode JSON response: %v", err)
		}

	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set(ContentType, ContentTypeJSON)
	w.WriteHeader(status)
	response := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}
