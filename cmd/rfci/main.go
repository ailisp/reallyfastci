package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.CORS())
	e.POST("/github", webhook.GithubWebhook)
	e.Static("/", "./public")
	e.Static("/build", "build")
	e.GET("/ws", notification.WebSocket)
	e.Logger.Fatal(e.Start(":1323"))
}
