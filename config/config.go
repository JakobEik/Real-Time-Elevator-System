package config

import (
	"fmt"
	"time"
)

const N_FLOORS = 4
const N_BUTTONS = 3
const N_ELEVATORS = 3
const DoorOpenDuration = time.Second * 3
const ToEveryone = -1

var ElevatorID = 0

type Behavior int

const (
	IDLE Behavior = iota
	DOOR_OPEN
	MOVING
	UNAVAILABLE
)

func (b Behavior) String() string {
	switch b {
	case IDLE:
		return "IDLE"
	case DOOR_OPEN:
		return "Door open"
	case MOVING:
		return "MOVING"
	default:
		return fmt.Sprintf("%d", int(b))
	}
}

type MessageType int

const (
	NEW_ORDER MessageType = iota
	DO_ORDER
	UPDATE_GLOBAL_STATE
	MSG_RECEIVED
	LOCAL_STATE_CHANGED
	GLOBAL_HALL_ORDERS
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
