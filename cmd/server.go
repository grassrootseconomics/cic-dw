package main

import (
	"cic-dw/internal/admin"
	"cic-dw/internal/dashboard"
	"cic-dw/internal/public"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initHTTPServer() *echo.Echo {
	server := echo.New()
	server.HideBanner = true
	// TODO: Remove after stable release
	server.Debug = true
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     conf.Server.Cors,
		AllowCredentials: true,
		MaxAge:           600,
	}))

	dashboard.InitDashboardApi(server, db, preparedQueries.dashboard)
	public.InitPublicApi(server, db, batchBalance, cicnetClient, preparedQueries.public)
	admin.InitAdminApi(server, db, preparedQueries.admin, metaClient, conf.Jwt.Secret)

	return server
}
