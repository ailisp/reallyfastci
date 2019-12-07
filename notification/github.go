package notification

import (
	"fmt"
	"log"

	"github.com/ailisp/reallyfastci/config"
	"github.com/ailisp/reallyfastci/core"
	"github.com/parnurzeal/gorequest"
)

func NotifyBuildStatusGithub(event *core.BuildEvent) {
	switch event.Status {
	case core.BuildSucceed:
		notifyGithub(event, "success")
	case core.BuildFailed:
		notifyGithub(event, "failure")
	case core.BuildQueued:
		notifyGithub(event, "pending")
	case core.BuildCancelled:
		notifyGithub(event, "failure")
	}
}

func notifyGithub(event *core.BuildEvent, status string) {
	request := gorequest.New()
	githubUrl := fmt.Sprintf("https://api.github.com/repos/%v/statuses/%v", config.RepoName, event.Commit)
	rfciUrl := fmt.Sprintf("%v/#build/%v", config.Config.ReallyfastciUrl, event.Commit)
	log.Printf("githubUrl: %v", githubUrl)
	log.Printf("rfciUrl: %v", rfciUrl)
	resp, body, errs := request.Post(githubUrl).
		Set("Authorization", fmt.Sprintf("token %v", config.Config.GithubToken)).
		Set("Accept", "application/vnd.github.antiope-preview+json").
		Send(fmt.Sprintf(`{"state":"%v","target_url":"%v","context":"reallyfastci"}`, status, rfciUrl)).
		End()
	if len(errs) > 0 {
		log.Printf("Error updating github status: %+v", errs)
	} else {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("Successfully updating github status")
		} else {
			log.Printf("Error updating github status: %+v", body)
		}
	}
}
