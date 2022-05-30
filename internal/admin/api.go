package admin

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nleof/goyesql"
	"github.com/rs/zerolog/log"
)

type api struct {
	db     *pgxpool.Pool
	q      goyesql.Queries
	jwtKey []byte
}

func InitAdminApi(e *echo.Echo, db *pgxpool.Pool, queries goyesql.Queries, jwtKey string) {
	api := newApi(db, queries, jwtKey)

	auth := e.Group(("/auth"))
	g := e.Group("/admin")

	auth.Use(api.dwCoreMiddleware)
	auth.POST("/login", sendLoginJwtCookie)
	auth.POST("/logout", sendLogoutCookie)

	g.Use(api.dwCoreMiddleware)
	g.Use(api.verifyAuthMiddleware)

	g.GET("/protected", handleProtectedResource)
}

func newApi(db *pgxpool.Pool, queries goyesql.Queries, jwtKey string) *api {
	return &api{
		db:     db,
		q:      queries,
		jwtKey: []byte(jwtKey),
	}
}

func (a *api) dwCoreMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("api", &api{
			db:     a.db,
			q:      a.q,
			jwtKey: a.jwtKey,
		})
		return next(c)
	}
}

func (a *api) verifyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Info().Msgf("%v", c.Cookies())
		cookie, err := c.Cookie("_ge_auth")
		if err != nil {
			return c.String(http.StatusForbidden, "auth cookie missing")
		}

		claims := &jwtClaims{}

		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return a.jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.String(http.StatusUnauthorized, "jwt signature validation failed")
			}
			return c.String(http.StatusBadRequest, "jwt bad request")
		}
		if !token.Valid {
			return c.String(http.StatusUnauthorized, "jwt invalid")
		}

		return next(c)
	}
}
