package config

import "Project/elevator"

const N_FLOORS = 4
const N_BUTTONS = 3
const N_ELEVATORS = 3
const DoorOpenDuration = 3

type Behavior int

const (
	Idle        Behavior = 0
	DoorOpen             = 1
	Moving               = 2
	Unavailable          = 3
)

type GlobalState struct {
	ElevatorID int
	states     []elevator.ElevatorState
}

/*type RequestState int

const (
	None      RequestState = 0
	NewOrder  RequestState = 1
	Confirmed RequestState = 2
	Complete  RequestState = 3
)*/
