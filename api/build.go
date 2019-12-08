package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ailisp/reallyfastci/build"
	"github.com/ailisp/reallyfastci/core"
	"github.com/hpcloud/tail"
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

func RunningOutput(c echo.Context) (err error) {
	commit := c.Param("commit")
	exitCodeFilename := fmt.Sprintf("build/%v/exitcode", commit)
	outputFilename := fmt.Sprintf("build/%v/output_combined.log", commit)
	if _, err := os.Stat(outputFilename); err == nil {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		outputTail, _ := tail.TailFile(outputFilename, tail.Config{
			Follow: true, ReOpen: true, MustExist: false,
			Logger: tail.DiscardingLogger,
		})
		exitcodeTail, _ := tail.TailFile(exitCodeFilename, tail.Config{
			Logger:    tail.DiscardingLogger,
			MustExist: false, ReOpen: true, Follow: true,
		})
		for {
			select {
			case line := <-outputTail.Lines:
				io.WriteString(c.Response(), line.Text)
				io.WriteString(c.Response(), "\n")
				c.Response().Flush()
			case <-exitcodeTail.Lines:
				return nil
			}
		}
	} else {
		return c.NoContent(http.StatusNotFound)
	}
}
