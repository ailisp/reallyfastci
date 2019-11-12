package core

import "fmt"

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
		fmt.Printf("%v", event)
	}
}
