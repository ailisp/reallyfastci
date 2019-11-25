package machine

import (
	"fmt"
	"sync"
	"time"

	"github.com/ailisp/reallyfastci/config"
	"github.com/cornelk/hashmap"
	"github.com/google/uuid"
)

type machineManager struct {
	runningMachines *hashmap.HashMap
	idleMachines    chan chan *Machine
	machineRequests chan *MachineRequest
	stopChan        chan bool
	maxMachines     int
	maxIdleMachines int

	numMachines           int
	numMachineMutex       *sync.Mutex
	machineDeletedSignals chan bool
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
	manager.idleMachines = make(chan chan *Machine, int(config.Config.Machine.IdleMachines))
	manager.maxMachines = int(config.Config.Machine.MaxMachines)
	manager.maxIdleMachines = int(config.Config.Machine.IdleMachines)
	manager.numMachines = 0
	manager.numMachineMutex = &sync.Mutex{}
	manager.machineDeletedSignals = make(chan bool, manager.maxMachines)
	go runMachineManager()
}

func runMachineManager() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case _ = <-manager.stopChan:
			break
		case req := <-manager.machineRequests:
			handleReq(req)
		case <-ticker.C:
			if manager.maxIdleMachines > 0 {
				startIdleMachines()
			}
		case <-manager.machineDeletedSignals:
			manager.numMachineMutex.Lock()
			manager.numMachines--
			manager.numMachineMutex.Unlock()
		}
	}
}

func handleReq(req *MachineRequest) {
	switch req.Op {
	case opRequestMachine:
		if _, ok := manager.runningMachines.GetStringKey(req.RequestId.String()); ok {
			req.RequestMachineChan <- nil
			return
		}

		manager.runningMachines.Set(req.RequestId.String(), &Machine{})

		if manager.maxIdleMachines > 0 {
			go handleRequestMachine(req.RequestId, req.RequestMachineChan)
		} else {
			go handleNoIdleRequestMachine(req.RequestId, req.RequestMachineChan)
		}
	case opReleaseMachine:
		go handleDeleteMachine(req.RequestId)
	}
}

func handleRequestMachine(requestId uuid.UUID, machineChan chan *Machine) {
	machine := <-<-manager.idleMachines
	if machine != nil {
		manager.runningMachines.Set(requestId.String(), machine)
	} else {
		manager.runningMachines.Del(requestId.String())
	}
	machineChan <- machine
}

func handleNoIdleRequestMachine(requestId uuid.UUID, machineChan chan *Machine) {
	machine := <-manager.newMachine()
	if machine != nil {
		manager.runningMachines.Set(requestId.String(), machine)
	} else {
		manager.runningMachines.Del(requestId.String())
	}
	machineChan <- machine
}

func (manager *machineManager) newMachine() (machineChan chan *Machine) {
	machineChan = manager.tryNewMachine()
	if machineChan != nil {
		return
	}

	tick := time.NewTicker(5 * time.Second)

	for {
		_ = <-tick.C
		machineChan = manager.tryNewMachine()
		if machineChan != nil {
			return
		}
	}
}

func (manager *machineManager) tryNewMachine() (machineChan chan *Machine) {
	manager.numMachineMutex.Lock()
	if manager.numMachines < manager.maxMachines {
		machineChan = newMachine(fmt.Sprintf("%v-%v", config.Config.Machine.Prefix, uuid.New().String()), manager.machineDeletedSignals)
		manager.numMachines++
	}
	manager.numMachineMutex.Unlock()

	return
}

func (manager *machineManager) deleteMachine(machine *Machine) (err error) {
	err = machine.delete()
	if err == nil {
		manager.machineDeletedSignals <- true
	}
	return
}

func handleDeleteMachine(requestId uuid.UUID) {
	val, ok := manager.runningMachines.GetStringKey(requestId.String())
	if ok {
		machine := val.(*Machine)
		if machine.Name != "" {
			manager.deleteMachine(machine)
			manager.runningMachines.Del(requestId.String())
		} else {
			ReleaseMachine(requestId)
		}
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

func startIdleMachines() {
	idleMachines := len(manager.idleMachines)
	for i := 0; i < manager.maxIdleMachines-idleMachines; i++ {
		machineChan := manager.tryNewMachine()
		if machineChan != nil {
			manager.idleMachines <- machineChan
		}
	}
}
