package distributor

import (
	drv "Project/driver"
	e "Project/elevator"
	"Project/network/peers"
	util "Project/utilities"
	//"Project/network/peers"
)

func Distributor(
	ch_doOrder chan<- drv.ButtonEvent,
	ch_localStateUpdated <-chan e.ElevatorState,
	ch_newLocalOrder <-chan drv.ButtonEvent,
	ch_peerUpdate chan peers.PeerUpdate,
	ch_peerTxEnable chan bool,
	ch_messageToNetwork chan<- util.NetworkMessage,
	ch_messageFromNetwork <-chan util.NetworkMessage) {
	var masterID = 0

	//globalState := util.InitGlobalState()
	localElevatorState := e.InitElev()
	// Ask network if they have a global state: true => globalState = this state, else
	if false {
		//TODO
	}
	for {
		select {
		case newLocalOrder := <-ch_newLocalOrder:
			newLocalOrderEvent(newLocalOrder, ch_messageToNetwork, masterID)
			ch_doOrder <- newLocalOrder
		case newLocalState := <-ch_localStateUpdated:
			localStateUpdatedEvent(newLocalState, localElevatorState, ch_messageToNetwork, masterID)
		case msg := <-ch_messageFromNetwork:
			newMessageEvent(msg, ch_doOrder)
		}

	}
}

func newLocalOrderEvent(order drv.ButtonEvent, ch_messageToNetwork chan<- util.NetworkMessage, masterID int) {

}

func localStateUpdatedEvent(
	newState e.ElevatorState,
	oldState e.ElevatorState,
	ch_messageToNetwork chan<- util.NetworkMessage,
	masterID int) {

	oldState = newState
	msg := util.CreateMessage(masterID, masterID, newState, util.LocalStateChange)
	ch_messageToNetwork <- msg

}
func newMessageEvent(msg util.NetworkMessage, ch_doOrder chan<- drv.ButtonEvent) {
	switch msg.MsgType {
	case util.DoOrder:
		m := msg.Msg.(drv.ButtonEvent)
		ch_doOrder <- m
	}
}
