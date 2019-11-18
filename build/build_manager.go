package build

import (
	"github.com/ailisp/reallyfastci/core"
	"github.com/cornelk/hashmap"
)

type buildManager struct {
	pendingPrBuilds   chan *core.PrEvent
	pendingPushBuilds chan *core.PushEvent
	runningBuilds     *hashmap.HashMap
	buildFinishEvents chan *core.BuildEvent
}

var manager buildManager

func InitBuildManager() {
	initBuildEventAgent()
	manager = buildManager{
		pendingPrBuilds:   make(chan *core.PrEvent, 100),
		pendingPushBuilds: make(chan *core.PushEvent, 100),
		buildFinishEvents: make(chan *core.BuildEvent, 100),
		runningBuilds:     &hashmap.HashMap{},
	}

	go runBuildManager()
}

func runBuildManager() {
	for {
		select {
		case push := <-manager.pendingPushBuilds:
			runPushBuild(push)
		case pr := <-manager.pendingPrBuilds:
			runPrBuild(pr)
		case buildFinish := <-manager.buildFinishEvents:
			_, ok := manager.runningBuilds.GetStringKey(buildFinish.Commit)
			if ok {
				manager.runningBuilds.Del(buildFinish.Commit)
			}
		}
	}
}

func QueuePushBuild(push *core.PushEvent) {
	manager.pendingPushBuilds <- push
}

func QueuePrBuild(pr *core.PrEvent) {
	manager.pendingPrBuilds <- pr
}

func runPushBuild(push *core.PushEvent) {
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

func runPrBuild(pr *core.PrEvent) {
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
