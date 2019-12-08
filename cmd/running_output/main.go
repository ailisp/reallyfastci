package main

import (
	"github.com/ailisp/reallyfastci/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/api/build/:commit/output", api.RunningOutput)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Info("shutting down the server")
	}
}
