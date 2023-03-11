package config

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

type NetworkMessage struct {
	SenderID   int
	ReceiverID int
	Msg        any
}

/*type RequestState int

const (
	None      RequestState = 0
	NewOrder  RequestState = 1
	Confirmed RequestState = 2
	Complete  RequestState = 3
)*/
