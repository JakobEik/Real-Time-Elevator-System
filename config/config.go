package config

const N_FLOORS = 4
const N_BUTTONS = 3
const DoorOpenDuration = 3

type RequestState int

const (
	None      RequestState = 0
	NewOrder  RequestState = 1
	Confirmed RequestState = 2
	Complete  RequestState = 3
)

type Behavior int

const (
	Idle        Behavior = 0
	DoorOpen    Behavior = 1
	Moving      Behavior = 2
	Unavailable Behavior = 3
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
