package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ailisp/reallyfastci/api"
	"github.com/ailisp/reallyfastci/build"
	"github.com/ailisp/reallyfastci/config"
	"github.com/ailisp/reallyfastci/machine"
	"github.com/ailisp/reallyfastci/notification"
	"github.com/ailisp/reallyfastci/webhook"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	config.LoadConfig("./config.yaml")

	machine.InitMachineManager()
	build.InitBuildManager()
	webhook.InitWebhook()
	notification.InitSse()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.CORS())
	e.POST("/github", webhook.GithubWebhook)
	e.GET("/api/build/:commit/output", api.RunningOutput)
	e.GET("/api/build/:commit", api.Build)
	e.GET("/api/build/:commit/exitcode", api.BuildExitCode)
	e.GET("/api/build", api.Index)
	e.GET("/sse", func(c echo.Context) (err error) {
		notification.SseServer.HTTPHandler(c.Response().Writer, c.Request())
		return nil
	})
	e.Static("/", ".")

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%v", config.Config.ApiPort)); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
