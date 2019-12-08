package main

import (
	"net/http"

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
	e.GET("/api/build/:commit", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "Running",
		})
	})
	if err := e.Start(":1323"); err != nil {
		e.Logger.Info("shutting down the server")
	}
}
