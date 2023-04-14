package config

import (
	"fmt"
	"time"
)

const N_FLOORS = 7
const N_BUTTONS = 3
const N_ELEVATORS = 3

const ConfirmationWaitDuration = time.Millisecond * 30

var ElevatorID int
var MasterID = 0

type Behavior int

const (
	IDLE Behavior = iota
	DOOR_OPEN
	MOVING
)

func (b Behavior) String() string {
	switch b {
	case IDLE:
		return "IDLE"
	case DOOR_OPEN:
		return "DOOR_OPEN"
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
	HALL_LIGHTS_UPDATE
)

func (t MessageType) String() string {
	switch t {
	case NEW_ORDER:
		return "NEW_ORDER"
	case DO_ORDER:
		return "DO_ORDER"
	case UPDATE_GLOBAL_STATE:
		return "UPDATE_GLOBAL_STATE"
	case MSG_RECEIVED:
		return "MSG_RECEIVED"
	case LOCAL_STATE_CHANGED:
		return "LOCAL_STATE_CHANGED"
	case HALL_LIGHTS_UPDATE:
		return "HALL_LIGHTS_UPDATE"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

type NetworkMessage struct {
	SenderID   int
	ReceiverID int
	Type       MessageType
	Content    any
}
