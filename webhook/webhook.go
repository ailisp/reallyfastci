package webhook

import (
	"log"
	"net/http"

	"github.com/ailisp/reallyfastci/core"
	"github.com/labstack/echo/v4"
)

func GithubWebhook(c echo.Context) (err error) {
	header := c.Request().Header

	eventType := header.Get("x-github-event")
	switch eventType {
	case "pull_request":
		pr := new(core.PrEvent)
		if err = c.Bind(pr); err != nil {
			log.Printf("%v", err)
			return c.JSON(http.StatusBadRequest, "")
		}
		if err = c.Validate(pr); err != nil {
			log.Printf("%v", err)
			return c.JSON(http.StatusBadRequest, "")
		}
		prs <- pr
		return c.JSON(http.StatusOK, "")
	case "push":
		push := new(core.PushEvent)
		if err = c.Bind(push); err != nil {
			log.Printf("%v", err)
			return c.JSON(http.StatusBadRequest, "")
		}
		if err = c.Validate(push); err != nil {
			log.Printf("%v", err)
			return c.JSON(http.StatusBadRequest, "")
		}
		pushes <- push
		return c.JSON(http.StatusOK, "")
	default:
		return c.JSON(http.StatusOK, "Ignored")
	}
}

func InitWebhook() {
	initPrAgent()
	initPushAgent()
}
