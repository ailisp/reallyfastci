package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ailisp/reallyfastci/build"
	"github.com/labstack/echo/v4"
)

type builds struct {
	Running  []string `json:"running"`
	Finished []string `json:"finished"`
}

func Index(c echo.Context) (err error) {
	outputBuilds, err := ioutil.ReadDir("./build")
	if err != nil {
		log.Println(err)
	}
	runningBuilds := build.RunningBuilds()

	finishedBuilds := []string{}

	for _, dir := range outputBuilds {
		buildCommit := dir.Name()
		_, ok := runningBuilds.GetStringKey(buildCommit)
		if !ok {
			finishedBuilds = append(finishedBuilds, buildCommit)
		}
	}

	runningBuildsArray := []string{}
	for build := range runningBuilds.Iter() {
		runningBuildsArray = append(runningBuildsArray, build.Key.(string))
	}

	return c.JSON(http.StatusOK, &builds{
		Running:  runningBuildsArray,
		Finished: finishedBuilds,
	})
}
