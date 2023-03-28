package main

import (
	"Project/assigner"
	"Project/config"
	d "Project/distributor"
	drv "Project/driver"
	e "Project/elevator"
	"Project/network/bcast"
	"Project/network/peers"
	"os"
	"strconv"
)

const bufferSize = 512

func main() {

	port := os.Args[1]
	config.ElevatorID, _ = strconv.Atoi(os.Args[2])

	// channels for Network
	ch_peerUpdate := make(chan peers.PeerUpdate)
	ch_peerTxEnable := make(chan bool)
	ch_messageToNetwork := make(chan config.Packet, bufferSize)
	ch_messageFromNetwork := make(chan config.Packet, bufferSize)

	// channels between distributor and FSM
	ch_localStateUpdated := make(chan e.ElevatorState, bufferSize)
	ch_executeOrder := make(chan drv.ButtonEvent, bufferSize)
	ch_globalHallOrders := make(chan [][]bool, bufferSize)

	// Channels for driver
	ch_buttonPress := make(chan drv.ButtonEvent, bufferSize)
	ch_floorArrival := make(chan int, bufferSize)
	ch_obstruction := make(chan bool, bufferSize)
	ch_stop := make(chan bool)

	//drv.Init("localhost:15657", config.N_FLOORS)
	drv.Init("localhost:"+port, config.N_FLOORS)

	// Driver go routines
	go drv.PollButtons(ch_buttonPress)
	go drv.PollFloorSensor(ch_floorArrival)
	go drv.PollObstructionSwitch(ch_obstruction)
	go drv.PollStopButton(ch_stop)

	// Networking go routines
	go bcast.Transmitter(20321, ch_messageToNetwork)
	go bcast.Receiver(20321, ch_messageFromNetwork)

	go d.Distributor(
		ch_messageFromNetwork,
		ch_messageToNetwork,
		ch_executeOrder,
		ch_globalHallOrders,
		ch_localStateUpdated,
		ch_buttonPress,
	)

	go assigner.MasterNode(
		ch_peerUpdate,
		ch_peerTxEnable,
		ch_messageFromNetwork,
		ch_messageToNetwork)

	e.Fsm(ch_executeOrder, ch_floorArrival, ch_obstruction, ch_stop, ch_localStateUpdated, ch_globalHallOrders)
}
