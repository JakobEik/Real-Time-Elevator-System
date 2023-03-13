package distributor

import (
	drv "Project/driver"
	e "Project/elevator"
	"Project/network/peers"
	"Project/utilities"
	//"Project/network/peers"
)

func Distributor(
	ch_doOrder chan<- drv.ButtonEvent,
	ch_localStateUpdated <-chan e.ElevatorState,
	ch_newLocalOrder <-chan drv.ButtonEvent,
	ch_peerUpdate chan peers.PeerUpdate,
	ch_peerTxEnable chan bool,
	ch_messageToNetwork chan<- utilities.NetworkMessage,
	ch_messageFromNetwork <-chan utilities.NetworkMessage) {

	globalState := initGlobalState()
	localElevatorState := e.InitElev()
	// Ask network if they have a global state: true => globalState = this state, else
	if false {
		//TODO
	}
	for {
		select {
		case newLocalOrder := <-ch_newLocalOrder:
			newLocalOrderEvent(newLocalOrder, ch_messageToNetwork)
			ch_doOrder <- newLocalOrder
		case newLocalState := <-ch_localStateUpdated:
			localElevatorState = newLocalState

		}

	}
}

func newLocalOrderEvent(order drv.ButtonEvent, ch_messageToNetwork chan<- utilities.NetworkMessage) {

}

func localStateUpdatedEvent(
	newState e.ElevatorState,
	oldState e.ElevatorState,
	ch_messageToNetwork chan<- utilities.NetworkMessage) {

	oldState = newState

}
