package machine

import (
	"fmt"

	"github.com/ailisp/reallyfastci/config"
	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
)

type machineManager struct {
	runningMachines *hashmap.HashMap
	idleMachines    chan *Machine
	machineRequests chan *MachineRequest
	stopChan        chan bool
}

var manager machineManager

type MachineRequest struct {
	Op int

	RequestId uuid.UUID

	RequestMachineChan chan *Machine
}

const (
	opRequestMachine = iota
	opReleaseMachine
)

func InitMachineManager() {
	manager.runningMachines = &hashmap.HashMap{}
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
		default:
			examineMachines()
		}
	}
}

func handleReq(req *MachineRequest) {
	switch req.Op {
	case opRequestMachine:
		if _, ok := manager.runningMachines.Get(req.RequestId); ok {
			return
		}
		manager.runningMachines.Set(req.RequestId, &Machine{})
		go handleRequestMachine(req.RequestId, req.RequestMachineChan)
	case opReleaseMachine:
		go handleDeleteMachine(req.RequestId)
	}
}

func handleRequestMachine(requestId uuid.UUID, machineChan chan *Machine) {
	machine := <-manager.idleMachines
	if machine != nil {
		manager.runningMachines.Set(requestId, machine)
	}
	machineChan <- machine
}

func handleDeleteMachine(requestId uuid.UUID) {
	machine, ok := manager.runningMachines.Get(requestId)
	if ok {
		machine.(*Machine).delete()
		manager.runningMachines.Del(requestId)
	}
}

func RequestMachine(requestId uuid.UUID) (machineChan chan *Machine) {
	machineChan = make(chan *Machine)
	manager.machineRequests <- &MachineRequest{Op: opRequestMachine,
		RequestId:          requestId,
		RequestMachineChan: machineChan,
	}
	return
}

func ReleaseMachine(requestId uuid.UUID) {
	manager.machineRequests <- &MachineRequest{Op: opReleaseMachine,
		RequestId: requestId,
	}
}

func examineMachines() {
	idleMachines := uint64(len(manager.idleMachines))
	runningMachines := uint64(manager.runningMachines.Len())
	if idleMachines < config.Config.Machine.IdleMachines && idleMachines+runningMachines < config.Config.Machine.MaxMachines {
		manager.idleMachines <- newMachine(fmt.Sprintf("%v-%v", config.Config.Machine.Prefix, uuid.New().String()))
	}
}
