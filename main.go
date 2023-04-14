package main

import (
	"Project/assigner"
	"Project/config"
	d "Project/distributor"
	drv "Project/driver"
	e "Project/elevator"
	"Project/network/bcast"
	"Project/network/peers"
	"Project/watchdog"
	"os"
	"strconv"
)

const bufferSize = 512

func main() {

	port := os.Args[1]
	config.ElevatorID, _ = strconv.Atoi(os.Args[2])
	ElevatorStrID := os.Args[2]

	// channels for Network
	ch_peerUpdate := make(chan peers.PeerUpdate)
	ch_peerTxEnable := make(chan bool)
	ch_packetToNetwork := make(chan d.Packet, bufferSize)
	ch_packetFromNetwork := make(chan d.Packet, bufferSize)

	// channels between distributor and FSM
	ch_localStateUpdated := make(chan e.ElevatorState, bufferSize)
	ch_executeOrder := make(chan drv.ButtonEvent, bufferSize)
	ch_globalHallOrders := make(chan [][]bool, bufferSize)

	// Channels for Packet Distributor
	ch_msgToPack := make(chan config.NetworkMessage, bufferSize)
	ch_msgToAssigner := make(chan config.NetworkMessage, bufferSize)
	ch_msgToDistributor := make(chan config.NetworkMessage, bufferSize)

	// Channels for driver
	ch_buttonPress := make(chan drv.ButtonEvent, bufferSize)
	ch_floorArrival := make(chan int, bufferSize)
	ch_obstruction := make(chan bool, bufferSize)
	ch_stop := make(chan bool)

	// Channels for watchdog
	ch_wdstart := make(chan bool)
	ch_wdstop := make(chan bool)
	ch_watchdogStuckBark := make(chan bool)

	//drv.Init("localhost:15657", config.N_FLOORS)
	drv.Init("localhost:"+port, config.N_FLOORS)

	// Driver go routines
	go drv.PollButtons(ch_buttonPress)
	go drv.PollFloorSensor(ch_floorArrival)
	go drv.PollObstructionSwitch(ch_obstruction)
	go drv.PollStopButton(ch_stop)

	// Networking go routines
	go bcast.Transmitter(23456, ch_packetToNetwork)
	go bcast.Receiver(23456, ch_packetFromNetwork)
	go peers.Transmitter(34567, ElevatorStrID, ch_peerTxEnable)
	go peers.Receiver(34567, ch_peerUpdate)
	
	go watchdog.Watchdog(ch_wdstart, ch_wdstop, ch_watchdogStuckBark)

	go d.Distributor(
		ch_msgToDistributor,
		ch_msgToPack,
		ch_executeOrder,
		ch_globalHallOrders,
		ch_localStateUpdated,
		ch_buttonPress,
		ch_watchdogStuckBark,
	)

	go d.PacketDistributor(
		ch_packetFromNetwork,
		ch_packetToNetwork,
		ch_msgToPack,
		ch_msgToAssigner,
		ch_msgToDistributor)

	go assigner.MasterNode(
		ch_peerUpdate,
		ch_msgToAssigner,
		ch_msgToPack)

	e.Fsm(
		ch_executeOrder,
		ch_floorArrival,
		ch_obstruction,
		ch_stop,
		ch_localStateUpdated,
		ch_globalHallOrders,
    	ch_wdstart, ch_wdstop,
		ch_peerTxEnable)
}
