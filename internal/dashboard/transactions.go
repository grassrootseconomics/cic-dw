package dashboard

import (
	"context"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
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

func handleLatestTokenTransactions(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		token = c.Param("address")

		data []tokenTransactionsRes
	)

	rows, err := api.db.Query(context.Background(), api.q["latest-token-transactions"], token)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	if len(data) < 1 {
		data = []tokenTransactionsRes{}
	}

	return c.JSON(http.StatusOK, data)
}
