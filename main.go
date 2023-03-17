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
	ch_messageToNetwork := make(chan config.NetworkMessage, bufferSize*10)
	ch_messageFromNetwork := make(chan config.NetworkMessage, bufferSize*10)

	// channels for distributor
	ch_localStateUpdated := make(chan e.ElevatorState)
	ch_newLocalOrder := make(chan drv.ButtonEvent)

	// channels for FSM
	ch_doOrder := make(chan drv.ButtonEvent, bufferSize)
	ch_floorArrival := make(chan int, bufferSize)
	ch_obstruction := make(chan bool, bufferSize)
	ch_stop := make(chan bool)

	// Channels for Elevio driver
	//ch_buttons := make(chan drv.ButtonEvent, bufferSize)

	//channel_DoorTimer := make(chan bool)

	drv.Init("localhost:15657", config.N_FLOORS)
	//drv.Init("localhost:12346", config.N_FLOORS)

	// Driver go routines
	go drv.PollButtons(ch_newLocalOrder)
	go drv.PollFloorSensor(ch_floorArrival)
	go drv.PollObstructionSwitch(ch_obstruction)
	go drv.PollStopButton(ch_stop)

	// Networking go routines
	go bcast.Transmitter(20321, ch_messageToNetwork)
	go bcast.Receiver(20321, ch_messageFromNetwork)

	go d.Distributor(
		ch_doOrder,
		ch_localStateUpdated,
		ch_newLocalOrder,
		ch_peerUpdate,
		ch_peerTxEnable,
		ch_messageToNetwork,
		ch_messageFromNetwork)

	e.Fsm(ch_doOrder, ch_floorArrival, ch_obstruction, ch_stop, ch_localStateUpdated)
}
