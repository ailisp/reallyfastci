package machine

import (
	"github.com/ailisp/reallyfastci/script"
	"github.com/google/uuid"
)

type Machine struct {
	name string
}

func newMachine() (machine *Machine, err error) {
	machine = &Machine{
		name: uuid.New().String(),
	}

	err = script.Run("pipenv", "run", "python", "create_machine.py")
	return machine, err
}

func (machine *Machine) delete() {
	script.DeleteMachine(machine.name)
}

func (machine *Machine) CloneRepo(url string, branch string, commit string) {
	script.CloneRepo(machine.name, url, branch, commit)
}

func (machine *Machine) CopyBuildScript(buildScript string) {
	script.CopyBuildScript(machine.name, buildScript)
}

func (machine *Machine) RunBuild(buildScript string) (exitCode int) {
	return script.RunBuild(machine.name, buildScript)
}
