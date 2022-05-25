package public

import (
	batch_balance "github.com/grassrootseconomics/cic-go/batch_balance"
	cic_net "github.com/grassrootseconomics/cic-go/net"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nleof/goyesql"
)

type api struct {
	db *pgxpool.Pool
	q  goyesql.Queries
	bb *batch_balance.BatchBalance
	cn *cic_net.CicNet
}

func InitPublicApi(e *echo.Echo, db *pgxpool.Pool, batchBalance *batch_balance.BatchBalance, cicnet *cic_net.CicNet, queries goyesql.Queries) {
	g := e.Group("/public")

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("api", &api{
				db: db,
				q:  queries,
				cn: cicnet,
				bb: batchBalance,
			})
			return next(c)
		}
	})

	// TODO: paginate schema validation

	g.GET("/balances/:address", handleBalancesQuery)
	g.GET("/tokens-count", handleTokensCountQuery)
	g.GET("/tokens", handleTokenListQuery)
	g.GET("/token/:address", handleTokenInfo)
	g.GET("/token-summary/:address", handleTokenSummary)
	g.GET("/latest-token-transactions/:address", handleTokenTransactions)
}
