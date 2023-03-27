package config

import "time"

const N_FLOORS = 4
const N_BUTTONS = 3
const N_ELEVATORS = 3
const DoorOpenDuration = time.Second * 3
const ToEveryone = -1

var ElevatorID = 0

type Behavior int

const (
	Idle Behavior = iota
	DoorOpen
	Moving
	Unavailable
)

type MessageType int

const (
	NewOrder MessageType = iota
	DoOrder
	UpdateGlobalState
	MsgReceived
	LocalStateChange
)

type NetworkMessage struct {
	SenderID   int
	ReceiverID int
	MasterID   int
	MsgType    MessageType
	Content    any
}

type Packet struct {
	Msg      NetworkMessage
	Checksum int
}
