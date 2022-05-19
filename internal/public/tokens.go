package public

import (
	"cic-dw/pkg/pagination"
	"context"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
)

type tokensRes struct {
	Id           int    `db:"id" json:"id"`
	TokenSymbol  string `db:"token_symbol" json:"token_symbol"`
	TokenName    string `db:"token_name" json:"token_name"`
	TokenAddress string `db:"token_address" json:"token_addres"`
}

type tokenCountRes struct {
	Count int `db:"count" json:"count"`
}

func handleTokenListQuery(c echo.Context) error {
	var (
		api = c.Get("api").(*api)
		pg  = pagination.GetPagination(c.QueryParams())
		res []tokensRes
		q   string
	)

	if pg.Forward {
		q = api.q["list-tokens-fwd"]
	} else {
		q = api.q["list-tokens-bkwd"]
	}

	rows, err := api.db.Query(context.Background(), q, pg.Cursor, pg.PerPage)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&res, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func handleTokensCountQuery(c echo.Context) error {
	var (
		api = c.Get("api").(*api)
		res tokenCountRes
	)

	rows, err := api.db.Query(context.Background(), api.q["tokens-count"])
	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(&res, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
