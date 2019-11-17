package core

import "github.com/ailisp/reallyfastci/machine"

type Build struct {
	repo        string
	branch      string
	commit      string
	buildScript string
	status      int

	eventChan *chan *BuildEvent
	cancel    chan bool
	machine   *machine.Machine
}

type BuildEvent struct {
	commit string
	status int
}

// BuildStatus
const (
	BuildQueued int = iota
	BuildMachineStarted
	BuildRepoCloned
	BuildScriptCopied

	BuildSucceed
	BuildFailed
	BuildCancelled
)

type newBuildParam struct {
	repo        string
	branch      string
	commit      string
	buildScript string
	eventChan   *chan *BuildEvent
}

func newBuild(param *newBuildParam) *Build {
	build := &Build{
		repo:        param.repo,
		branch:      param.branch,
		commit:      param.commit,
		buildScript: param.buildScript,

		eventChan: param.eventChan,
		cancel:    make(chan bool),
	}
	build.updateStatus(BuildQueued)

	go build.run()
	return build
}

func (build *Build) run() {
	build.machine = <-machine.RequestCreateMachine()
	build.updateStatus(BuildMachineStarted)

	build.machine.CloneRepo(build.repo, build.branch, build.commit)
	build.updateStatus(BuildRepoCloned)

	build.machine.CopyBuildScript()
	build.updateStatus(BuildScriptCopied)

	err := build.machine.RunBuild(build.commit)
	if err == nil {
		build.updateStatus(BuildSucceed)
	} else {
		build.updateStatus(BuildFailed)
	}

}

func (build Build) updateStatus(status int) {
	build.status = status
	(*build.eventChan) <- &BuildEvent{
		commit: build.commit,
		status: build.status,
	}
}

func (build Build) sendCancel() {
	build.cancel <- true
}
