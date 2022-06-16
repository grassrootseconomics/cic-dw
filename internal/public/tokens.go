package public

import (
	"cic-dw/pkg/address"
	"context"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
	"github.com/lmittmann/w3"
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

type TokenInfoRes struct {
	IsDemurrage bool   `json:"is_demurrage"`
	Name        string `json:"token_name"`
	Symbol      string `json:"token_symbol"`
	TotalSupply int64  `json:"token_total_supply"`
}

type tokenSummaryRes struct {
	TotalHolders      int64 `db:"count" json:"token_holders"`
	TotalTransactions int64 `db:"count" json:"token_transactions"`
}

func handleTokenListQuery(c echo.Context) error {
	var (
		api = c.Get("api").(*api)

		res []tokensRes
	)

	rows, err := api.db.Query(context.Background(), api.q["list-tokens"])
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

func handleTokenInfo(c echo.Context) error {
	var (
		api          = c.Get("api").(*api)
		tokenAddress = c.Param("address")
		rCtx         = context.Background()

		res TokenInfoRes
	)

	_, err := api.cn.DemurrageTokenInfo(rCtx, w3.A(address.Checksum(tokenAddress)))
	if err != nil {
		res.IsDemurrage = false
	} else {
		res.IsDemurrage = true
	}

	tokenInfo, err := api.cn.ERC20TokenInfo(rCtx, w3.A(address.Checksum(tokenAddress)))
	if err != nil {
		return err
	}

	res.Name = tokenInfo.Name
	res.Symbol = tokenInfo.Symbol
	res.TotalSupply = tokenInfo.TotalSupply.Int64() / 1000000

	return c.JSON(http.StatusOK, res)
}

func handleTokenSummary(c echo.Context) error {
	var (
		api   = c.Get("api").(*api)
		token = c.Param("address")

		data tokenSummaryRes
	)

	uniqueTokenHoldersrRow, err := api.db.Query(context.Background(), api.q["unique-token-holders"], token)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(&data.TotalHolders, uniqueTokenHoldersrRow); err != nil {
		return err
	}

	tokenTxRow, err := api.db.Query(context.Background(), api.q["all-time-token-transactions-count"], token)
	if err != nil {
		return err
	}

	if err := pgxscan.ScanOne(&data.TotalTransactions, tokenTxRow); err != nil {

		return err
	}

	return c.JSON(http.StatusOK, data)
}
