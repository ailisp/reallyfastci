package machine

import (
	"github.com/ailisp/reallyfastci/config"
	"github.com/ailisp/reallyfastci/script"
	"github.com/google/uuid"
	"strconv"
)

type Machine struct {
	name string
}

func newMachine() (machine *Machine) {
	machine = &Machine{
		name: uuid.New().String(),
	}

	err := script.Run("create_machine.py",
		"--name", machine.name,
		"--machine_type", config.Config.Machine.MachineType,
		"--disk_size", strconv.FormatUint(config.Config.Machine.DiskSizeGB, 10),
		"--image_project", config.Config.Machine.ImageProject,
		"--image_family", config.Config.Machine.ImageFamily,
		"--zone", config.Config.Machine.Zone)
	if err != nil {
		return &Machine{}
	} else {
		return machine
	}
}

func (machine *Machine) delete() error {
	return script.Run("delete_machine.py",
		"--name", machine.name)
}

func (machine *Machine) CloneRepo(url string, branch string, commit string) error {
	return script.Run("clone_repo_on_machine.py",
		"--name", machine.name,
		"--url", url,
		"--branch", branch,
		"--commit", commit)
}

func (machine *Machine) CopyBuildScript() error {
	return script.Run("copy_build_script_to_machine.py",
		"--name", machine.name,
		"--local_path", config.Config.Build.Script)
}

func (machine *Machine) RunBuild(commit string) (err error) {
	return script.Run("run_build.py",
		"--name", machine.name,
		"--commit", commit,
		"--local_path", config.Config.Build.Script)
}
