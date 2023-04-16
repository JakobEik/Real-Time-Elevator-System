package config

import (
	"Project/driver"
	"fmt"
)

type ElevatorState struct {
	Floor     int
	Direction driver.MotorDirection
	Behavior  Behavior
	Orders    [][]bool
}

type Behavior int

const (
	EB_IDLE Behavior = iota
	EB_DOOR_OPEN
	EB_MOVING
)

type MessageType int

const (
	NEW_ORDER MessageType = iota
	DO_ORDER
	UPDATE_GLOBAL_STATE
	MSG_RECEIVED
	LOCAL_STATE_CHANGED
	HALL_LIGHTS_UPDATE
	NEW_MASTER
)

type NetworkMessage struct {
	SenderID   int
	ReceiverID int
	Type       MessageType
	Content    any
}

// String methods

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
	case NEW_MASTER:
		return "NEW ORDER"
	default:
		return fmt.Sprintf("%d", int(t))
	}
}

func (b Behavior) String() string {
	switch b {
	case EB_IDLE:
		return "EB_IDLE"
	case EB_DOOR_OPEN:
		return "EB_DOOR_OPEN"
	case EB_MOVING:
		return "EB_MOVING"
	default:
		return fmt.Sprintf("%d", int(b))
	}
}
