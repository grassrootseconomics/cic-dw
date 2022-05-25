package dashboard

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nleof/goyesql"
)

type api struct {
	db *pgxpool.Pool
	q  goyesql.Queries
}

func InitDashboardApi(e *echo.Echo, db *pgxpool.Pool, queries goyesql.Queries) {
	g := e.Group("/dashboard")

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("api", &api{
				db: db,
				q:  queries,
			})
			return next(c)
		}
	})

	g.GET("/new-registrations", handleNewRegistrations)
	g.GET("/transactions-count", handleTransactionsCount)
	g.GET("/token-transactions-count/:address", handleTokenTransactionsCount)
	g.GET("/token-volume/:address", handleTokenVolume)
	g.GET("/latest-token-transactions/:address", handleLatestTokenTransactions)
}
