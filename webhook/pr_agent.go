package webhook

import (
	"fmt"

	"github.com/ailisp/reallyfastci/build"
	"github.com/ailisp/reallyfastci/core"
)

var prs chan *core.PrEvent

func initPrAgent() {
	prs = make(chan *core.PrEvent, 100)

	go runPrAgent()
}

func runPrAgent() {
	for {
		event := <-prs
		fmt.Printf("Received a pull request event: %+v\n", event)
		if event.Action != "closed" {
			build.QueuePrBuild(event)
		}
	}
}
