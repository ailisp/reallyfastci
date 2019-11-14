package core

import "github.com/cornelk/hashmap"

type buildManager struct {
	pendingPrBuilds   chan *PrEvent
	pendingPushBuilds chan *PushEvent
	runningBuilds     *hashmap.HashMap
	buildFinishEvents chan *BuildEvent
}

var BuildManager buildManager

func InitBuildManager() {
	BuildManager = buildManager{
		pendingPrBuilds:   make(chan *PrEvent, 100),
		pendingPushBuilds: make(chan *PushEvent, 100),
		buildFinishEvents: make(chan *BuildEvent, 100),
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
		case buildFinish := <-manager.buildFinishEvents:
			_, ok := manager.runningBuilds.GetStringKey(buildFinish.commit)
			if ok {
				manager.runningBuilds.Del(buildFinish.commit)
			}
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
		prevBuild.(*Build).sendCancel()
	}

	prevCommitBuild, ok := manager.runningBuilds.GetStringKey(push.Before)
	if ok {
		prevCommitBuild.(*Build).sendCancel()
	}

	manager.runningBuilds.Set(push.After, newBuild(&newBuildParam{
		repo:   push.Repo.HtmlUrl,
		branch: push.Ref,
		commit: push.After,
	}))
}

func (manager buildManager) runPrBuild(pr *PrEvent) {
	if pr.Before != "" {
		prevBuild, ok := manager.runningBuilds.GetStringKey(pr.Before)
		if ok {
			prevBuild.(*Build).sendCancel()
		}
	}
	if pr.After != "" {
		prevBuild, ok := manager.runningBuilds.GetStringKey(pr.After)
		if ok {
			prevBuild.(*Build).sendCancel()
		}
	}
	manager.runningBuilds.Set(pr.PullRequest.Head.Sha, &newBuildParam{
		repo:   pr.PullRequest.Head.Repo.HtmlUrl,
		branch: pr.PullRequest.Head.Ref,
		commit: pr.PullRequest.Head.Sha,
	})
}
