package core

import "github.com/cornelk/hashmap"

type buildManager struct {
	pendingPrBuilds   chan *PrEvent
	pendingPushBuilds chan *PushEvent
	runningBuilds     *hashmap.HashMap
}

var BuildManager buildManager

func InitBuildManager() {
	BuildManager = buildManager{
		pendingPrBuilds:   make(chan *PrEvent, 100),
		pendingPushBuilds: make(chan *PushEvent, 100),
		runningBuilds:     &hashmap.HashMap{},
	}

	go BuildManager.run()
}

func (manager buildManager) run() {
	for {
		select {
		case push := <-manager.pendingPushBuilds:
			manager.runPushBuild(push)
		case pr := <-manager.pendingPrBuilds:
			manager.runPrBuild(pr)
		}
	}
}

func (manager buildManager) queuePushBuild(push *PushEvent) {
	manager.pendingPushBuilds <- push
}

func (manager buildManager) queuePrBuild(pr *PrEvent) {
	manager.pendingPrBuilds <- pr
}

func (manager buildManager) runPushBuild(push *PushEvent) {
	prevBuild, ok := manager.runningBuilds.GetStringKey(push.After)
	if ok {
		prevBuild.(Build).sendStop()
	}

	prevCommitBuild, ok := manager.runningBuilds.GetStringKey(push.Before)
	if ok {
		prevCommitBuild.(Build).sendStop()
	}

	manager.runningBuilds.Set(push.After, newBuild(push.Repo.HtmlUrl, push.Ref, push.After))
}

func (manager buildManager) runPrBuild(pr *PrEvent) {
	if pr.Before != "" {
		prevBuild, ok := manager.runningBuilds.GetStringKey(pr.Before)
		if ok {
			prevBuild.(Build).sendStop()
		}
	}
	if pr.After != "" {
		prevBuild, ok := manager.runningBuilds.GetStringKey(pr.After)
		if ok {
			prevBuild.(Build).sendStop()
		}
	}
	manager.runningBuilds.Set(pr.PullRequest.Head.Sha, newBuild(pr.PullRequest.Head.Repo.HtmlUrl, pr.PullRequest.Head.Ref, pr.PullRequest.Head.Sha))
}
