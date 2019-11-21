package build

import (
	"github.com/ailisp/reallyfastci/core"
	"github.com/ailisp/reallyfastci/machine"
)

type Build struct {
	repo        string
	branch      string
	commit      string
	buildScript string
	status      int

	cancel  chan bool
	machine *machine.Machine
}

type newBuildParam struct {
	repo        string
	branch      string
	commit      string
	buildScript string
	eventChan   chan *core.BuildEvent
}

func newBuild(param *newBuildParam) *Build {
	build := &Build{
		repo:        param.repo,
		branch:      param.branch,
		commit:      param.commit,
		buildScript: param.buildScript,

		cancel: make(chan bool),
	}
	build.updateStatus(core.BuildQueued)

	go build.run()
	return build
}

func (build *Build) run() {
	for {
		switch build.status {
		case core.BuildQueued:
			machineName, machineChan := machine.RequestCreateMachine()
			select {
			case <-build.cancel:
				machine.RequestDeleteMachine(machineName)
				build.updateStatus(core.BuildCancelled)
				return
			case m := <-machineChan:
				if m != nil {
					build.machine = m
					build.updateStatus(core.BuildMachineStarted)
				} else {
					build.updateStatus(core.BuildFailed)
					return
				}
			}
		case core.BuildMachineStarted:
			errChan := build.machine.CloneRepo(build.repo, build.branch, build.commit)
			select {
			case <-build.cancel:
				machine.RequestDeleteMachine(build.machine.Name)
				build.updateStatus(core.BuildCancelled)
				return
			case err := <-errChan:
				if err == nil {
					build.updateStatus(core.BuildRepoCloned)
				} else {
					build.updateStatus(core.BuildFailed)
					return
				}
			}
		case core.BuildRepoCloned:
			errChan := build.machine.CopyBuildScript()
			select {
			case <-build.cancel:
				machine.RequestDeleteMachine(build.machine.Name)
				build.updateStatus(core.BuildCancelled)
				return
			case err := <-errChan:
				if err == nil {
					build.updateStatus(core.BuildScriptCopied)
				} else {
					build.updateStatus(core.BuildFailed)
					return
				}
			}
		case core.BuildScriptCopied:
			errChan := build.machine.RunBuild(build.commit)
			select {
			case <-build.cancel:
				machine.RequestDeleteMachine(build.machine.Name)
				build.updateStatus(core.BuildCancelled)
				return
			case err := <-errChan:
				if err == nil {
					build.updateStatus(core.BuildSucceed)
					return
				} else {
					build.updateStatus(core.BuildFailed)
					return
				}
			}
		}
	}
}

func (build *Build) updateStatus(status int) {
	build.status = status
	buildEventChan <- &core.BuildEvent{
		Commit: build.commit,
		Status: build.status,
	}
}

func (build *Build) sendCancel() {
	build.cancel <- true
}
