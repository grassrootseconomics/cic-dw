package pagination

import (
	"net/url"
	"strconv"
)

type Pagination struct {
	PerPage   int
	Cursor    int
	Next      bool
	FirstPage bool
}

func GetPagination(q url.Values) Pagination {
	var (
		pp, _     = strconv.Atoi(q.Get("per_page"))
		cursor, _ = strconv.Atoi(q.Get("cursor"))
		next, _   = strconv.ParseBool(q.Get("next"))
		firstPage = false
	)

	if pp > 100 {
		pp = 100
	}

	if !next && cursor < 1 {
		firstPage = true
	}

	return Pagination{
		PerPage:   pp,
		Cursor:    cursor,
		Next:      next,
		FirstPage: firstPage,
	}
}
