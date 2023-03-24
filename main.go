package main

import (
	"Project/assigner"
	"Project/config"
	d "Project/distributor"
	drv "Project/driver"
	e "Project/elevator"
	"Project/network"
	"Project/network/bcast"
	"Project/network/peers"
	"os"
	"strconv"
)

const bufferSize = config.N_ELEVATORS * 11

func main() {

	port := os.Args[1]
	config.ElevatorID, _ = strconv.Atoi(os.Args[2])

	// channels for Network
	ch_peerUpdate := make(chan peers.PeerUpdate)
	ch_peerTxEnable := make(chan bool)
	ch_messageToNetwork := make(chan config.Packet, 512)
	ch_messageFromNetwork := make(chan config.Packet, 512)

	// Channels between master and network
	ch_calculateNewOrder := make(chan drv.ButtonEvent, 32)
	ch_localStateToMaster := make(chan assigner.StateMessage, 32)
	ch_orderCalculated := make(chan assigner.OrderMessage, 32)
	ch_updateGlobalStateDemand := make(chan []e.ElevatorState, 32)

	// Channels between distributor and network
	ch_doOrder := make(chan drv.ButtonEvent, 32)
	ch_updateGlobalState := make(chan []e.ElevatorState, 32)
	ch_localStateFromLocal := make(chan e.ElevatorState, 32)
	ch_newLocalOrder := make(chan drv.ButtonEvent, 32)

	// channels between distributor and FSM
	ch_localStateUpdated := make(chan e.ElevatorState)
	ch_executeOrder := make(chan drv.ButtonEvent, bufferSize)

	// Channels between driver and FSM
	ch_buttonPress := make(chan drv.ButtonEvent)
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
		ch_doOrder,
		ch_localStateUpdated,
		ch_buttonPress,
		ch_updateGlobalState,
		ch_localStateFromLocal,
		ch_newLocalOrder,
		ch_executeOrder)

	go assigner.MasterNode(
		ch_peerUpdate,
		ch_peerTxEnable,
		ch_calculateNewOrder,
		ch_localStateToMaster,
		ch_orderCalculated,
		ch_updateGlobalStateDemand)

	go network.NetworkHandler(
		ch_messageFromNetwork, 
		ch_messageToNetwork, 
		ch_calculateNewOrder, 
		ch_localStateToMaster, 
		ch_doOrder,
		ch_updateGlobalState,
		ch_orderCalculated,
		ch_updateGlobalStateDemand, ch_localStateFromLocal, ch_newLocalOrder)

	e.Fsm(ch_executeOrder, ch_floorArrival, ch_obstruction, ch_stop, ch_localStateUpdated)
}
