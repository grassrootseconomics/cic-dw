package public

import (
	"cic-dw/pkg/pagination"
	"context"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type tokenTransactionsRes struct {
	Id      int64     `db:"id" json:"id"`
	Block   int64     `db:"block_number" json:"block"`
	Date    time.Time `db:"date_block" json:"time"`
	TxHash  string    `db:"tx_hash" json:"tx_hash"`
	From    string    `db:"sender_address" json:"from"`
	To      string    `db:"recipient_address" json:"to"`
	Value   int64     `db:"tx_value" json:"tx_value"`
	Success bool      `db:"success" json:"success"`
}

func handleTokenTransactions(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		token = c.Param("address")
		pg    = pagination.GetPagination(c.QueryParams())

		data []tokenTransactionsRes
	)

	if pg.Cursor == -1 {
		var max int64
		if err := api.db.QueryRow(context.Background(), "SELECT MAX(id) + 10 from transactions").Scan(&max); err != nil {
			return err
		}

		pg.Cursor = int(max)
	}

	log.Info().Msgf("%d", pg.Cursor)

	rows, err := api.db.Query(context.Background(), api.q["latest-token-transactions"], token, pg.Cursor, pg.PerPage)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}
