package build

import (
	"log"

	"github.com/ailisp/reallyfastci/core"
	"github.com/ailisp/reallyfastci/notification"
)

var buildEventChan chan *core.BuildEvent

func initBuildEventAgent() {
	buildEventChan = make(chan *core.BuildEvent, 100)
	go runBuildEventAgent()
}

func runBuildEventAgent() {
	for {
		buildEvent := <-buildEventChan
		log.Printf("Build event: %+v", buildEvent)
		notification.NotifyBuildStatusGithub(buildEvent)
		notification.NotifySse(buildEvent)
		if buildEvent.Status >= core.BuildSucceed {
			manager.buildFinishEvents <- buildEvent
		}
	}
}
