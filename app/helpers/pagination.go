package helpers

import (
	"net/http"
	"strconv"
)

// Pagination
func Pagination(r *http.Request, limit int) (int, int) {
	keys := r.URL.Query()
	if keys.Get("page") == "" {
		return 1, 0
	}

	page, _ := strconv.Atoi(keys.Get("page"))
	if page < 1 {
		return 1, 0
	}
	begin := (limit * page) - limit
	return page, begin
}