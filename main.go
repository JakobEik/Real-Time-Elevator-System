package main

import (
	"Project/config"
	d "Project/distributor"
	drv "Project/driver"
	e "Project/elevator"
	"Project/network/bcast"
	"Project/network/peers"
)

const bufferSize = config.N_ELEVATORS * 11

func main() {
	// channels for Network
	ch_peerUpdate := make(chan peers.PeerUpdate)
	ch_peerTxEnable := make(chan bool)
	// Transmitter channels
	ch_TxGlobalState := make(chan config.NetworkMessage, bufferSize)
	ch_TxNewOrder := make(chan config.NetworkMessage, bufferSize)
	ch_TxOrderDone := make(chan config.NetworkMessage, bufferSize)
	ch_TxOrderAccepted := make(chan config.NetworkMessage, bufferSize)
	ch_TxRequestGlobalState := make(chan config.NetworkMessage, config.N_ELEVATORS+1)
	ch_TxDoOrder := make(chan config.NetworkMessage, bufferSize)
	ch_TxChangeYourState := make(chan config.NetworkMessage)
	ch_TxMsgReceived := make(chan bool)

	// Receiver channels
	ch_RxGlobalState := make(chan config.NetworkMessage, bufferSize)
	ch_RxNewOrder := make(chan config.NetworkMessage, bufferSize)
	ch_RxOrderDone := make(chan config.NetworkMessage, bufferSize)
	ch_RxOrderAccepted := make(chan config.NetworkMessage, bufferSize)
	ch_RxRequestGlobalState := make(chan config.NetworkMessage, config.N_ELEVATORS+1)
	ch_RxDoOrder := make(chan config.NetworkMessage, bufferSize)
	ch_RxChangeYourState := make(chan config.NetworkMessage, bufferSize)
	ch_RxMsgReceived := make(chan config.NetworkMessage, bufferSize)

	// channels for distributor
	ch_localStateUpdated := make(chan e.ElevatorState)
	ch_newLocalOrder := make(chan drv.ButtonEvent)

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
	go bcast.Transmitter(15600,
		ch_TxGlobalState,
		ch_TxNewOrder,
		ch_TxOrderDone,
		ch_TxOrderAccepted,
		ch_TxRequestGlobalState,
		ch_TxDoOrder,
		ch_TxChangeYourState,
		ch_TxMsgReceived)
	go bcast.Receiver(15600,
		ch_RxGlobalState,
		ch_RxNewOrder,
		ch_RxOrderDone,
		ch_RxOrderAccepted,
		ch_RxRequestGlobalState,
		ch_RxDoOrder,
		ch_RxChangeYourState,
		ch_RxMsgReceived)

	go d.Distributor(ch_doOrder, ch_localStateUpdated, ch_newLocalOrder, ch_peerUpdate, ch_peerTxEnable, ch_TxGlobalState, ch_RxGlobalStats)
	e.Fsm(ch_buttons, ch_doOrder, ch_floorArrival, ch_obstruction, ch_stop)
}
