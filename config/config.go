package config

import "Project/driver"

const N_FLOORS = 4
const N_BUTTONS = 3
const DoorOpenDuration = 3

type DirectionBehaviourPair struct {
	Direction driver.MotorDirection
	Behaviour Behaviour
}

type RequestState int

const (
	None      RequestState = 0
	NewOrder  RequestState = 1
	Confirmed RequestState = 2
	Complete  RequestState = 3
)

type Behaviour int

const (
	Idle        Behaviour = 0
	DoorOpen    Behaviour = 1
	Moving      Behaviour = 2
	Unavailable Behaviour = 3
)

type ButtonType int

const (
	HallUp   ButtonType = 0
	HallDown ButtonType = 1
	Cab      ButtonType = 2
)

type Order struct {
	Floor  int
	Button ButtonType
}

type CostRequest struct {
	Id         string
	Cost       int
	AssignedID string
	Req        Order
}
