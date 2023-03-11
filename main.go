package main

import (
	"Project/config"
	drv "Project/driver"
	e "Project/elevator"
	//"Project/network/bcast"
	//"Project/network/peers"
)

const bufferSize = config.N_ELEVATORS * 11

func main() {
	// channels for Network
	/*ch_peerUpdate := make(chan peers.PeerUpdate)
	ch_peerTxEnable := make(chan bool)
	// Transmitter channels
	ch_TxGlobalState := make(chan config.GlobalState, bufferSize)
	ch_TxNewOrder := make(chan drv.ButtonEvent, bufferSize)
	ch_TxOrderDone := make(chan drv.ButtonEvent, bufferSize)
	ch_TxOrderAccepted := make(chan drv.ButtonEvent, bufferSize)
	ch_TxRequestGlobalState := make(chan bool, config.N_ELEVATORS+1)
	ch_TxDoOrder := make(drv.ButtonEvent, bufferSize)
	ch_TxChangeYourState := make(TODO)
	ch_TxMsgReceived := make(chan bool)

	// Receiver channels
	ch_RxGlobalStats := make(chan config.GlobalState, bufferSize)
	ch_RxNewOrder := make(chan drv.ButtonEvent, bufferSize)
	ch_RxOrderDone := make(chan drv.ButtonEvent, bufferSize)
	ch_RxOrderAccepted := make(chan drv.ButtonEvent, bufferSize)
	ch_RxRequestGlobalState := make(chan bool, config.N_ELEVATORS+1)
	ch_RxDoOrder := make(drv.ButtonEvent, bufferSize)
	ch_RxChangeYourState := make(TODO)
	ch_RxMsgReceived := make(chan bool)*/

	// channels for distributor
	//ch_localStateUpdated := make(chan e.ElevatorState)
	//ch_newLocalOrder := make(chan drv.ButtonEvent)

	// channels for FSM
	ch_doOrder := make(chan drv.ButtonEvent, 50)
	ch_floorArrival := make(chan int, 50)
	ch_obstruction := make(chan bool, 50)
	ch_stop := make(chan bool)

	// Channels for Elevio driver
	ch_buttons := make(chan drv.ButtonEvent, 100)

	//channel_DoorTimer := make(chan bool)

	drv.Init("localhost:15657", config.N_FLOORS)
	// Driver go routines
	go drv.PollButtons(ch_buttons)
	go drv.PollFloorSensor(ch_floorArrival)
	go drv.PollObstructionSwitch(ch_obstruction)
	go drv.PollStopButton(ch_stop)

	// Networking go routines
	//go bcast.Transmitter()

	//go d.Distributor(ch_doOrder, ch_localStateUpdated, ch_newLocalOrder)
	e.Fsm(ch_buttons, ch_doOrder, ch_floorArrival, ch_obstruction, ch_stop)
}
