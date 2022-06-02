package admin

import (
	"context"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
)

type pinStatusResponse struct {
	PhoneNumber       string `db:"phone_number" json:"phone_number"`
	FailedPinAttempts int    `db:"failed_pin_attempts" json:"failed_pin_attempts"`
	AccountStatus     string `db:"account_status" json:"account_status"`
}

func handlePinStatus(c echo.Context) error {
	var (
		api = c.Get("api").(*api)
		res []pinStatusResponse
	)

	rows, err := api.db.Query(context.Background(), api.q["pin-status"])
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&res, rows); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
