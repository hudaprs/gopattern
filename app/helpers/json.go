package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSON returns a well formated response with a statusCode
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)

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
