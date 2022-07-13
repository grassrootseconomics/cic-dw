package admin

import (
	"cic-dw/pkg/pagination"
	"context"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type userTransactionRes struct {
	Id                  int64     `db:"id" json:"id"`
	Date                time.Time `db:"date_block" json:"tx_date"`
	TxHash              string    `db:"tx_hash" json:"tx_hash"`
	TokenSymbol         string    `db:"token_symbol" json:"voucher"`
	SenderAddress       string    `db:"sender_address" json:"sender_address"`
	RecipeintAddress    string    `db:"recipient_address" json:"recipient_address"`
	TxValue             int64     `db:"tx_value" json:"tx_value"`
	TxSuccess           bool      `db:"success" json:"tx_success"`
	SenderIdentifier    string    `db:"sender_identifier" json:"sender_identifier"`
	RecipientIdentifier string    `db:"recipient_identifier" json:"recipient_identifier"`
}

func handleLatestTransactions(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		phone = c.Param("phone")
		pg    = pagination.GetPagination(c.QueryParams())

		data []userTransactionRes
		rows pgx.Rows
		err  error
	)

	if pg.FirstPage {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions"], phone, pg.PerPage)
		if err != nil {
			return err
		}
	} else if pg.Next {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-next"], phone, pg.Cursor, pg.PerPage)
		if err != nil {
			return err
		}
	} else {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-previous"], phone, pg.Cursor, pg.PerPage)
		if err != nil {
			return err
		}
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func handleLatestTransactionsByToken(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		phone = c.Param("phone")
		token = c.Param("token")
		pg    = pagination.GetPagination(c.QueryParams())

		data []userTransactionRes
		rows pgx.Rows
		err  error
	)

	if pg.FirstPage {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-by-token"], phone, token, pg.PerPage)
		if err != nil {
			return err
		}
	} else if pg.Next {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-by-token-next"], phone, token, pg.Cursor, pg.PerPage)
		if err != nil {
			return err
		}
	} else {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-by-token-previous"], phone, token, pg.Cursor, pg.PerPage)
		if err != nil {
			return err
		}
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}

func handleLatestTransactionsByArchivedToken(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		phone = c.Param("phone")
		token = c.Param("token")
		pg    = pagination.GetPagination(c.QueryParams())

		data []userTransactionRes
		rows pgx.Rows
		err  error
	)

	if pg.FirstPage {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-by-archived-token"], phone, token, pg.PerPage)
		if err != nil {
			return err
		}
	} else if pg.Next {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-by-archived-token-next"], phone, token, pg.Cursor, pg.PerPage)
		if err != nil {
			return err
		}
	} else {
		rows, err = api.db.Query(context.Background(), api.q["account-latest-transactions-by-archived-token-previous"], phone, token, pg.Cursor, pg.PerPage)
		if err != nil {
			return err
		}
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, data)
}
