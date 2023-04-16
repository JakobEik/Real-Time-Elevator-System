package main

import (
	"Project/assigner"
	"Project/config"
	d "Project/distributor"
	drv "Project/driver"
	e "Project/elevator"
	"Project/failroutine"
	"Project/network/bcast"
	"Project/network/peers"
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

	// channels between distributor and Fsm
	ch_localStateUpdated := make(chan config.ElevatorState, bufferSize)
	ch_executeOrder := make(chan drv.ButtonEvent, bufferSize)
	ch_hallLights := make(chan [][]bool, bufferSize)
	ch_unavailable := make(chan bool, bufferSize)

	// Channels between Assigner and Fsm
	ch_offNetwork := make(chan bool, bufferSize)

	// Channels for Packet Distributor
	ch_packMessage := make(chan config.NetworkMessage, bufferSize)
	ch_msgToAssigner := make(chan config.NetworkMessage, bufferSize)
	ch_msgToDistributor := make(chan config.NetworkMessage, bufferSize)

	// Channels for driver
	ch_buttonPress := make(chan drv.ButtonEvent, bufferSize)
	ch_floorArrival := make(chan int, bufferSize)
	ch_obstruction := make(chan bool, bufferSize)
	ch_stop := make(chan bool)

	drv.Init("localhost:"+port, config.N_FLOORS)

	// Driver go routines
	go drv.PollButtons(ch_buttonPress)
	go drv.PollFloorSensor(ch_floorArrival)
	go drv.PollObstructionSwitch(ch_obstruction)
	go drv.PollStopButton(ch_stop)

	// Networking go routines
	go bcast.Transmitter(config.BROADCAST_PORT, ch_packetToNetwork)
	go bcast.Receiver(config.BROADCAST_PORT, ch_packetFromNetwork)
	go peers.Transmitter(config.PEER_PORT, ElevatorStrID, ch_peerTxEnable)
	go peers.Receiver(config.PEER_PORT, ch_peerUpdate)

	go d.Distributor(
		ch_msgToDistributor,
		ch_packMessage,
		ch_executeOrder,
		ch_hallLights,
		ch_localStateUpdated,
		ch_buttonPress,
		ch_unavailable,
		ch_peerTxEnable)

	go d.PacketDistributor(
		ch_packetFromNetwork,
		ch_packetToNetwork,
		ch_packMessage,
		ch_msgToAssigner,
		ch_msgToDistributor)

	go assigner.Assigner(
		ch_peerUpdate,
		ch_msgToAssigner,
		ch_packMessage,
		ch_offNetwork)

	e.Fsm(
		ch_executeOrder,
		ch_floorArrival,
		ch_obstruction,
		ch_stop,
		ch_localStateUpdated,
		ch_hallLights,
		ch_unavailable,
		ch_offNetwork)

	println("EXIT PROGRAM")
	failroutine.FailRoutine()

}
