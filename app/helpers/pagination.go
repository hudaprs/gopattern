package helpers

import (
	"net/http"
	"strconv"
)

/**
@desc PaginationResponse make pagination response way better

@param r *http.Request
@param limit int query string from param for limiting data
*/
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

/**
@desc PaginationResponse make pagination response way better

@param *http.Request
@param page int query string from param for getting data from different page
@param pages int list of total pages (total / limit)
@param limit int query string from param for limiting data
@param total int total data in database
@param data interface{} presenting the data there's gonna be output
*/
func PaginationResponse(r *http.Request, page, pages, limit, total int, data interface{}) interface{} {
	return map[string]interface{}{
		"Links": map[string]interface{}{
			"First": r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(page),
			"Last":  r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(pages),
			"Prev":  r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(page-1),
			"Next":  r.URL.Host + r.URL.Path + "?page=" + strconv.Itoa(page+1),
		},
		"Meta": map[string]interface{}{
			"Limit":        limit,
			"Total":        total,
			"TotalPage":    pages,
			"CurrentPage":  page,
			"NextPage":     page + 1,
			"PreviousPage": page - 1,
			"LastPage":     pages,
		},
		"Results": data,
	}
}
