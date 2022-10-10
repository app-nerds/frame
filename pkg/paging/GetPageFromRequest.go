package paging

import (
	"net/http"
	"strconv"
)

func GetPageFromRequest(r *http.Request) int {
	var (
		err        error
		pageString string
		page       int
	)

	if r.Method == http.MethodGet {
		pageString = r.URL.Query().Get("page")
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		pageString = r.FormValue("page")
	}

	if page, err = strconv.Atoi(pageString); err != nil {
		page = 1
	}

	return page
}
