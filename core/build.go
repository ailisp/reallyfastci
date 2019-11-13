package core

type Build struct {
	repo   string
	branch string
	commit string

	status chan int
	cancel chan bool
}

const (
	Start int = iota
	CreateMachine
	CopyRunner
	Clone
	Run
	Success
	Fail
	Cancel
	MachineDeleted
)

func newBuild(repo string, branch string, commit string) Build {
	build := Build{
		repo:   repo,
		branch: branch,
		commit: commit,

		status: make(chan int),
		cancel: make(chan bool),
	}
	build.status <- Start
	go build.run()
	return build
}

func (build Build) run() {
	for {

	}
}

func (build Build) sendStop() {
	build.cancel <- true
}
