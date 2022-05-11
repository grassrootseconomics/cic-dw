package dashboard

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type lineChartRes struct {
	X time.Time `db:"x" json:"x"`
	Y int       `db:"y" json:"y"`
}

func handleNewRegistrations(c echo.Context) error {
	var (
		api  = c.Get("api").(*api)
		data []lineChartRes
	)

	from, to := parseDateRange(c.QueryParams())

	rows, err := api.db.Query(context.Background(), api.q["new-user-registrations"], from, to)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func handleTransactionsCount(c echo.Context) error {
	var (
		api  = c.Get("api").(*api)
		data []lineChartRes
	)

	from, to := parseDateRange(c.QueryParams())

	rows, err := api.db.Query(context.Background(), api.q["transactions-count"], from, to)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}
