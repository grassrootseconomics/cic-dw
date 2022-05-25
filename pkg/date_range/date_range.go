package date_range

import (
	"net/url"

	"github.com/golang-module/carbon/v2"
)

func ParseDateRange(q url.Values) (string, string) {
	var from, to string

	qFrom := q.Get("from")
	qTo := q.Get("to")

	parseFrom := carbon.Parse(qFrom)
	parseTo := carbon.Parse(qTo)

	if parseFrom.Error != nil || parseTo.Error != nil || qFrom == "" || qTo == "" {
		from = carbon.Now().StartOfMonth().ToDateString()
		to = carbon.Now().EndOfMonth().ToDateString()
	} else {
		from = parseFrom.ToDateString()
		to = parseTo.ToDateString()
	}

	return from, to
}
