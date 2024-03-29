package dashboard

import (
	"cic-dw/pkg/date_range"
	"context"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type lineChartRes struct {
	X time.Time `db:"x" json:"x"`
	Y int       `db:"y" json:"y"`
}

func handleNewRegistrations(c echo.Context) error {
	var (
		api = c.Get("api").(*api)
		qP  = c.QueryParams()

		rows pgx.Rows
		err  error
		data []lineChartRes
	)

	from, to := date_range.ParseDateRange(qP)

	if qP.Get("country") == "cmr" {
		rows, err = api.db.Query(context.Background(), api.q["new-user-registrations-cmr"], from, to)
		if err != nil {
			return err
		}
	} else {
		rows, err = api.db.Query(context.Background(), api.q["new-user-registrations"], from, to)
		if err != nil {
			return err
		}
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func handleTransactionsCount(c echo.Context) error {
	var (
		api = c.Get("api").(*api)
		qP  = c.QueryParams()

		rows pgx.Rows
		err  error
		data []lineChartRes
	)

	from, to := date_range.ParseDateRange(qP)

	if qP.Get("country") == "cmr" {
		rows, err = api.db.Query(context.Background(), api.q["transactions-count-cmr"], from, to)
		if err != nil {
			return err
		}
	} else {
		rows, err = api.db.Query(context.Background(), api.q["transactions-count"], from, to)
		if err != nil {
			return err
		}
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func handleTokenTransactionsCount(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		token = c.Param("address")

		data []lineChartRes
	)

	from, to := date_range.ParseDateRange(c.QueryParams())

	rows, err := api.db.Query(context.Background(), api.q["token-transactions-count"], from, to, token)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func handleTokenVolume(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		token = c.Param("address")

		data []lineChartRes
	)

	from, to := date_range.ParseDateRange(c.QueryParams())

	rows, err := api.db.Query(context.Background(), api.q["token-volume"], from, to, token)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {

		return err
	}

	return c.JSON(http.StatusOK, data)
}
