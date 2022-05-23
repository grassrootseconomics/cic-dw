package public

import (
	"cic-dw/pkg/address"
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
	"github.com/lmittmann/w3"
)

type balanceRes struct {
	TokenSymbol string `json:"symbol"`
	Balance     int64  `json:"balance"`
}

type dbRes struct {
	TokenSymbol  string `db:"token_symbol"`
	TokenAddress string `db:"token_address"`
}

func handleBalancesQuery(c echo.Context) error {
	var (
		api            = c.Get("api").(*api)
		data           []dbRes
		tokenAddresses []common.Address
		res            []balanceRes
	)

	// TODO: return 400
	qAddress, err := address.SarafuAddress(c.Param("address"))
	if err != nil {
		return err
	}

	rows, err := api.db.Query(context.Background(), api.q["all-known-tokens"], qAddress)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanAll(&data, rows); err != nil {
		return err
	}

	for _, rowData := range data {
		tokenAddresses = append(tokenAddresses, w3.A(address.Checksum(rowData.TokenAddress)))
	}

	balances, err := api.bb.TokensBalance(context.Background(), w3.A(address.Checksum(qAddress)), tokenAddresses)
	if err != nil {
		return err
	}

	for i, balance := range balances {
		var d balanceRes

		d.Balance = balance.Int64()
		d.TokenSymbol = data[i].TokenSymbol

		res = append(res, d)
	}

	return c.JSON(http.StatusOK, res)
}
