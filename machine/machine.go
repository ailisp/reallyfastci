package machine

import "github.com/google/uuid" 

type Machine struct {
	name string
}

func newMachine() (machine *Machine) {
	machine = &Machine{
		name: uuid.New().String()
	}
	script.CreateMachine()
	return machine
}

func (machine *Machine) delete() {

}
