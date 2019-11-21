package webhook

import (
	"fmt"
	"strings"

	"github.com/ailisp/reallyfastci/build"
	"github.com/ailisp/reallyfastci/config"
	"github.com/ailisp/reallyfastci/core"
)

var pushes chan *core.PushEvent

func initPushAgent() {
	pushes = make(chan *core.PushEvent, 100)
	go runPushAgent()
}

func runPushAgent() {
	for {
		event := <-pushes
		fmt.Printf("Received a push event: %+v\n", event)
		ref := strings.Split(event.Ref, "/")
		if len(ref) > 0 {
			branch := ref[len(ref)-1]
			for _, build_branch := range config.Config.PushTriggerBranches {
				if build_branch == branch {
					build.QueuePushBuild(event)
				}
			}
		}
	}
}
