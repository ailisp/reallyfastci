package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ailisp/reallyfastci/build"
	"github.com/ailisp/reallyfastci/core"
	"github.com/labstack/echo/v4"
)

type buildStatus struct {
	Status string `json:"status"`
}

type buildOutput struct {
	Status         string `json:"status"`
	OutputCombined string `json:"output_combined"`
	ExitCode       int    `json:"exitcode"`
}

func Build(c echo.Context) (err error) {
	commit := c.Param("commit")
	exitCodeFilename := fmt.Sprintf("build/%v/exitcode", commit)
	if _, err := os.Stat(exitCodeFilename); err == nil {
		status := core.BuildStatusStr(core.BuildSucceed)
		exitCodeContent, _ := ioutil.ReadFile(exitCodeFilename)
		exitCode, _ := strconv.Atoi(string(exitCodeContent))
		if exitCode != 0 {
			status = core.BuildStatusStr(core.BuildFailed)
		}
		outputContent, _ := ioutil.ReadFile(fmt.Sprintf("build/%v/output_combined.log", commit))
		return c.JSON(http.StatusOK, &buildOutput{
			Status:         status,
			ExitCode:       exitCode,
			OutputCombined: string(outputContent),
		})
	} else {
		status := build.GetRunningBuild(commit)
		if status != core.BuildNotRunning {
			return c.JSON(http.StatusOK, &buildStatus{
				Status: core.BuildStatusStr(status),
			})
			// Frontend need to GET /ws to listen for build change
		} else {
			return c.JSON(http.StatusNotFound, nil)
		}
	}
}
