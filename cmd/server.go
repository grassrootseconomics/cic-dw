package main

import (
	"cic-dw/internal/dashboard"
	"github.com/labstack/echo/v4"
)

func initHTTPServer() *echo.Echo {
	server := echo.New()
	server.HideBanner = true

	// TODO: Remove after stable release
	server.Debug = true

	dashboard.InitDashboardApi(server, db, preparedQueries.dashboard)

	return server
}
