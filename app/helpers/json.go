package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Success response JSON
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"Status":  "Success",
		"Message": message,
		"Data": data,
	}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// Error response JSON
func Error(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if message == "" {
		message = "Something went wrong"
	}
	response := map[string]interface{}{
		"Status": "Error",
		"Message": message,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
