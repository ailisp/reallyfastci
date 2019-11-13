package core

import (
	"fmt"
	"strings"
)

type pushAgent struct {
	Pushs chan *PushEvent
}

var PushAgent pushAgent

func InitPushAgent() {
	PushAgent = pushAgent{Pushs: make(chan *PushEvent, 100)}
	go PushAgent.run()
}

func (agent pushAgent) Send(push *PushEvent) {
	agent.Pushs <- push
}

func (agent pushAgent) run() {
	for {
		event := <-agent.Pushs
		fmt.Printf("Received a push event: %+v\n", event)
		ref := strings.Split(event.Ref, "/")
		if len(ref) > 0 {
			branch := ref[len(ref)-1]
			if branch == "master" || branch == "staging" {
				BuildManager.queuePushBuild(event)
			}
		}
	}
}
