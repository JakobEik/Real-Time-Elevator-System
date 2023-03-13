package config

const N_FLOORS = 4
const N_BUTTONS = 3
const N_ELEVATORS = 3
const DoorOpenDuration = 3
const ToEveryone = -1
const ElevatorID = 0

type Behavior int

const (
	Idle        Behavior = 0
	DoorOpen             = 1
	Moving               = 2
	Unavailable          = 3
)

type MessageType int

const (
	GlobalState MessageType = iota
	NewOrder
	OrderDone
	OrderAccepted
	RequestGlobalState
	DoOrder
	ChangeYourState
	MsgReceived
	LocalStateChange
)

type NetworkMessage struct {
	SenderID   int
	ReceiverID int
	MasterID   int
	Msg        any
	MsgType    MessageType
}
