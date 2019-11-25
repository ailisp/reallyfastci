package script

import (
	"fmt"
	"log"
	"os/exec"
)

func Run(name string, arg ...string) (errChan chan error) {
	errChan = make(chan error)
	go func() {
		args := append([]string{"run", "python", fmt.Sprintf("script/%v", name)}, arg...)
		cmd := exec.Command("pipenv", args...)
		err := cmd.Run()
		if err != nil {
			log.Printf("script.Run() %v %v failed with %s\n", name, arg, err)
		}
		errChan <- err
	}()
	return errChan
}
