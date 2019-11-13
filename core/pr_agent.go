package core

import "fmt"

type prAgent struct {
	Prs chan *PrEvent
}

var PrAgent prAgent

func InitPrAgent() {
	PrAgent = prAgent{Prs: make(chan *PrEvent, 100)}
	go PrAgent.run()
}

func (agent prAgent) Send(pr *PrEvent) {
	agent.Prs <- pr
}

func (agent prAgent) run() {
	for {
		event := <-agent.Prs
		fmt.Printf("Received a pull request event: %+v\n", event)
		BuildManager.queuePrBuild(event)
	}
}
