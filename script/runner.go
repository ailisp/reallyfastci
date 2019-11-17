package script

import (
	"log"
	"os/exec"
)

func Run(name string, arg ...string) (err error) {
	args := append([]string{"run", "python", name}, arg...)
	cmd := exec.Command("pipenv", args...)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
