package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON returns a well formatted response with a statusCode
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

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

// ERROR returns a jsonified error response along with a status code
func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusInternalServerError, nil)
}
