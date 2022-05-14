package public

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nleof/goyesql"
)

type api struct {
	db *pgxpool.Pool
	q  goyesql.Queries
}

func InitPublicApi(e *echo.Echo, db *pgxpool.Pool, queries goyesql.Queries) {
	g := e.Group("/public")

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("api", &api{
				db: db,
				q:  queries,
			})
			return next(c)
		}
	})

	g.GET("/balances/:address", handleBalancesQuery)
}
