package distributor

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"Project/network/peers"
	"Project/utils"
	"fmt"
	//"Project/network/peers"
)

func Distributor(
	ch_doOrder chan<- drv.ButtonEvent,
	ch_localStateUpdated <-chan e.ElevatorState,
	ch_newLocalOrder <-chan drv.ButtonEvent,
	ch_peerUpdate chan peers.PeerUpdate,
	ch_peerTxEnable chan bool,
	ch_messageToNetwork chan<- c.NetworkMessage,
	ch_messageFromNetwork <-chan c.NetworkMessage) {
	var masterID = 0

	//globalState := c.InitGlobalState()
	localElevatorState := e.InitElev()
	// Ask network if they have a global state: true => globalState = this state, else
	if false {
		//TODO
	}
	for {
		select {
		case newLocalOrder := <-ch_newLocalOrder:
			newLocalOrderEvent(newLocalOrder, ch_messageToNetwork, masterID)
		case newLocalState := <-ch_localStateUpdated:
			localStateUpdatedEvent(newLocalState, localElevatorState, ch_messageToNetwork, masterID)
		case msg := <-ch_messageFromNetwork:
			newMessageEvent(msg, ch_doOrder)
		}

	}
}

func newLocalOrderEvent(order drv.ButtonEvent, ch_messageToNetwork chan<- c.NetworkMessage, masterID int) {
	msg := utils.CreateMessage(masterID, masterID, order, c.NewOrder)
	fmt.Println(msg)
	ch_messageToNetwork <- msg
	println("sent")
}

func localStateUpdatedEvent(
	newState e.ElevatorState,
	oldState e.ElevatorState, //pointer???
	ch_messageToNetwork chan<- c.NetworkMessage,
	masterID int) {

	oldState = newState
	msg := utils.CreateMessage(masterID, masterID, newState, c.LocalStateChange)
	ch_messageToNetwork <- msg

}
func newMessageEvent(msg c.NetworkMessage, ch_doOrder chan<- drv.ButtonEvent) {
	fmt.Println(msg)
	if msg.ReceiverID == c.ElevatorID {
		switch msg.MsgType {
		case c.DoOrder:
			m := msg.Content.(drv.ButtonEvent)
			ch_doOrder <- m

		}

	}
}
