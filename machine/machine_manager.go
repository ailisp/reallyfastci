package machine

import "github.com/cornelk/hashmap"

var machines *hashmap.HashMap
var machineRequests chan *MachineRequest
var stopChan chan bool

type MachineRequest struct {
	Op                int
	MachineChan       chan *Machine
	DeleteMachineName string
	DeleteFinishChan  chan bool
}

const (
	opCreateMachine = iota
	opDeleteMachine
)

func StartMachineManager() {
	machines = &hashmap.HashMap{}
	machineRequests = make(chan *MachineRequest, 100)
	stopChan = make(chan bool)
	go run()
}

func run() {
	for {
		select {
		case _ = <-stopChan:
			break
		case req := <-machineRequests:
			handleReq(req)
		}
	}
}

func handleReq(req *MachineRequest) {
	switch req.Op {
	case opCreateMachine:
		go createMachine(req.MachineChan)
	case opDeleteMachine:
		go deleteMachine(req.DeleteMachineName, req.DeleteFinishChan)
	}
}

func createMachine(machineChan chan *Machine) {
	machine := createMachine()
	machines.Set(machine.name, machine)
	machineChan <- machine
}

func deleteMachine(deleteMachineName string, deleteFinishChan chan bool) {
	machine, ok := machines.GetStringKey(deleteMachineName)
	if ok {
		machines.Del(deleteMachineName)
		machine.(*Machine).delete()
		deleteFinishChan <- true
	}
}

func RequestCreateMachine() (machineChan chan *Machine) {
	machineChan = make(chan *Machine)
	machineRequests <- &MachineRequest{Op: opCreateMachine,
		MachineChan: machineChan,
	}
	return machineChan
}

func RequestDeleteMachine(deleteMachineName string) (finishChan chan bool) {
	finishChan = make(chan bool)
	machineRequests <- &MachineRequest{Op: opDeleteMachine,
		DeleteFinishChan:  finishChan,
		DeleteMachineName: deleteMachineName,
	}
	return finishChan
}
