package machine

import "github.com/cornelk/hashmap"

type machineManager struct {
	machines        *hashmap.HashMap
	machineRequests chan *MachineRequest
	stopChan        chan bool
}

var manager machineManager

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

func InitMachineManager() {
	manager.machines = &hashmap.HashMap{}
	manager.machineRequests = make(chan *MachineRequest, 100)
	manager.stopChan = make(chan bool)
	go runMachineManager()
}

func runMachineManager() {
	for {
		select {
		case _ = <-manager.stopChan:
			break
		case req := <-manager.machineRequests:
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
	machine := newMachine()
	manager.machines.Set(machine.name, machine)
	machineChan <- machine
}

func deleteMachine(deleteMachineName string, deleteFinishChan chan bool) {
	machine, ok := manager.machines.GetStringKey(deleteMachineName)
	if ok {
		manager.machines.Del(deleteMachineName)
		machine.(*Machine).delete()
		deleteFinishChan <- true
	}
}

func RequestCreateMachine() (machineChan chan *Machine) {
	machineChan = make(chan *Machine)
	manager.machineRequests <- &MachineRequest{Op: opCreateMachine,
		MachineChan: machineChan,
	}
	return machineChan
}

func RequestDeleteMachine(deleteMachineName string) (finishChan chan bool) {
	finishChan = make(chan bool)
	manager.machineRequests <- &MachineRequest{Op: opDeleteMachine,
		DeleteFinishChan:  finishChan,
		DeleteMachineName: deleteMachineName,
	}
	return finishChan
}
