package core

type Build struct {
	repo   string
	branch string
	commit string
	status int

	eventChan chan *BuildEvent
	cancel    chan bool
}

type BuildEvent struct {
	commit string
	status int
}

// BuildStatus
const (
	BuildQueued int = iota
	BuildReadyToRun
	BuildRepoCloned

	BuildSucceed
	BuildFailed
	BuildCancelled
)

type newBuildParam struct {
	repo      string
	branch    string
	commit    string
	eventChan chan *BuildEvent
}

func newBuild(param *newBuildParam) *Build {
	build := &Build{
		repo:   param.repo,
		branch: param.branch,
		commit: param.commit,

		eventChan: newBuildParam.eventChan,
		cancel:    make(chan bool),
	}
	build.status = BuildQueued
	go build.run()
	return build
}

func (build Build) run() {
	for {
		switch build.status {
		case BuildQueued:
			machine := MachineManager.RequestMachine()
		}
	}
}

func (build Build) sendCancel() {
	build.cancel <- true
}
