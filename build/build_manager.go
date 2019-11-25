package build

import (
	"log"

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
			log.Printf("Build finsh event received: %+v, remove build from running builds", buildFinish)
			manager.runningBuilds.Del(buildFinish.Commit)
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
		log.Printf("Cancel build %+v due to a new push build with same commit comes", prevBuild)
		prevBuild.(*Build).sendCancel()
	}

	prevCommitBuild, ok := manager.runningBuilds.GetStringKey(push.Before)
	if ok {
		log.Printf("Cancel build %+v due to a new push to branch comes", prevCommitBuild)
		prevCommitBuild.(*Build).sendCancel()
	}

	log.Printf("Build scheduled to run for push %+v", push)
	manager.runningBuilds.Set(push.After, newBuild(&newBuildParam{
		repo:   push.Repo.HtmlUrl,
		branch: push.Branch,
		commit: push.After,
	}))
}

func runPrBuild(pr *core.PrEvent) {
	if pr.Before != "" {
		prevBuild, ok := manager.runningBuilds.GetStringKey(pr.Before)
		if ok {
			log.Printf("Cancel build %+v due to a new push to PR comes", prevBuild)
			prevBuild.(*Build).sendCancel()
		}
	}
	if pr.After != "" {
		prevBuild, ok := manager.runningBuilds.GetStringKey(pr.After)
		if ok {
			log.Printf("Cancel build %+v due to a new PR build with same commit hash comes", prevBuild)
			prevBuild.(*Build).sendCancel()
		}
	}

	log.Printf("Build scheduled to run for pr %+v", pr)
	manager.runningBuilds.Set(pr.PullRequest.Head.Sha, newBuild(&newBuildParam{
		repo:   pr.PullRequest.Head.Repo.HtmlUrl,
		branch: pr.PullRequest.Head.Ref,
		commit: pr.PullRequest.Head.Sha,
	}))
}
