package machine

import (
	"strconv"

	"github.com/ailisp/reallyfastci/config"
	"github.com/ailisp/reallyfastci/script"
)

type Machine struct {
	Name string
}

func newMachine(machineName string, failChan chan bool) (machineChan chan *Machine) {
	machine := &Machine{
		Name: machineName,
	}
	machineChan = make(chan *Machine)

	go func() {
		err := <-script.Run("create_machine.py",
			"--name", machine.Name,
			"--machine_type", config.Config.Machine.MachineType,
			"--disk_size", strconv.FormatUint(config.Config.Machine.DiskSizeGB, 10),
			"--image_project", config.Config.Machine.ImageProject,
			"--image_family", config.Config.Machine.ImageFamily,
			"--zone", config.Config.Machine.Zone)
		if err != nil {
			machineChan <- nil
			failChan <- true
		} else {
			machineChan <- machine
		}
	}()

	return
}

func (machine *Machine) delete() error {
	return <-script.Run("delete_machine.py",
		"--name", machine.Name)
}

func (machine *Machine) CloneRepo(url string, branch string, commit string) (errChan chan error) {
	return script.Run("clone_repo_on_machine.py",
		"--name", machine.Name,
		"--url", url,
		"--repo", config.RepoName,
		"--branch", branch,
		"--commit", commit)
}

func (machine *Machine) CopyBuildScript() (errChan chan error) {
	return script.Run("copy_build_script_to_machine.py",
		"--name", machine.Name,
		"--local_path", config.Config.Build.Script)
}

func (machine *Machine) RunBuild(commit string) (errChan chan error) {
	return script.Run("run_build.py",
		"--name", machine.Name,
		"--commit", commit,
		"--local_path", config.Config.Build.Script)
}
