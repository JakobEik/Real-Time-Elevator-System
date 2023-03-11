package distributor

import (
	"Project/config"
	"Project/driver"
	e "Project/elevator"
	//"Project/network/peers"
)

var globalState []e.ElevatorState

const elevatorID = 0

func Distributor(
	ch_doOrder chan<- driver.ButtonEvent,
	ch_localStateUpdated <-chan e.ElevatorState,
	ch_newLocalOrder <-chan driver.ButtonEvent,
	/*ch_peerUpdate chan peers.PeerUpdate,
	  ch_peerTxEnable chan bool,
	  ch_Tx chan []e.ElevatorState,
	  ch_Rx chan []e.ElevatorState*/) {

	// Ask network if they have a global state: true => globalState = this state, else
	if false {
		//TODO
	} else {
		initGlobalState()
	}
	for {
		select {
		case newLocalOrder := <-ch_newLocalOrder:
			//TODO: FIX NETWORKING STUFF
			ch_doOrder <- newLocalOrder
		case newLocalState := <-ch_localStateUpdated:

		}

	}
}

func initGlobalState() {
	globalState = make([]elevator.ElevatorState, config.N_ELEVATORS)
	for i := 0; i < config.N_ELEVATORS-1; i++ {
		globalState = append(globalState, elevator.InitElev())
	}
}
