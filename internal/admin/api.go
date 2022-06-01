package admin

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/grassrootseconomics/cic-go/meta"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nleof/goyesql"
	"github.com/rs/zerolog/log"
)

type api struct {
	db     *pgxpool.Pool
	q      goyesql.Queries
	m      *meta.CicMeta
	jwtKey []byte
}

func InitAdminApi(e *echo.Echo, db *pgxpool.Pool, queries goyesql.Queries, metaClient *meta.CicMeta, jwtKey string) {
	api := newApi(db, queries, metaClient, jwtKey)

	auth := e.Group(("/auth"))
	g := e.Group("/admin")

	auth.Use(api.dwCoreMiddleware)
	auth.POST("/login", sendLoginJwtCookie)
	auth.POST("/logout", sendLogoutCookie)

	g.Use(api.dwCoreMiddleware)
	g.Use(api.verifyAuthMiddleware)

	g.GET("/check", isLoggedIn)
	g.GET("/meta-proxy/:address", handleMetaProxy)
}

func newApi(db *pgxpool.Pool, queries goyesql.Queries, metaClient *meta.CicMeta, jwtKey string) *api {
	log.Info().Msgf("%s inj", jwtKey)
	return &api{
		db:     db,
		q:      queries,
		m:      metaClient,
		jwtKey: []byte(jwtKey),
	}
}

func (a *api) dwCoreMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("api", &api{
			db:     a.db,
			q:      a.q,
			m:      a.m,
			jwtKey: a.jwtKey,
		})
		return next(c)
	}
}

func (a *api) verifyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("_ge_auth")
		if err != nil {
			return c.String(http.StatusForbidden, "auth cookie missing")
		}

		claims := &jwtClaims{}

		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return a.jwtKey, nil
		})
		if err != nil {
			return c.String(http.StatusUnauthorized, "jwt validation failed")
		}
		if !token.Valid {
			return c.String(http.StatusUnauthorized, "jwt invalid")
		}

		return next(c)
	}
}
