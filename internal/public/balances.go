package public

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
	"net/http"
)

type dbRes struct {
	TokenSymbol  string `db:"token_symbol"`
	TokenAddress string `db:"token_address"`
}

func handleBalancesQuery(c echo.Context) error {
	var (
		api  = c.Get("api").(*api)
		data []dbRes
	)

	rows, err := api.db.Query(context.Background(), api.q["all-known-tokens"], c.Param("address"))
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}
