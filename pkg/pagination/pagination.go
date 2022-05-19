package pagination

import (
	"net/url"
	"strconv"
)

type Pagination struct {
	PerPage int
	Cursor  int
	Forward bool
}

func GetPagination(q url.Values) Pagination {
	var (
		pp, _      = strconv.Atoi(q.Get("per_page"))
		cursor, _  = strconv.Atoi(q.Get("cursor"))
		forward, _ = strconv.ParseBool(q.Get("forward"))
	)

	if pp > 100 {
		pp = 100
	}

	if !forward && cursor < 1 {
		cursor = 1
	}

	return Pagination{
		PerPage: pp,
		Cursor:  cursor,
		Forward: forward,
	}
}
