package utilities

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
)

type NetworkMessage struct {
	SenderID   int
	ReceiverID int
	Msg        any
	MsgType    MessageType
}

/*type RequestState int

const (
	None      RequestState = 0
	NewOrder  RequestState = 1
	Confirmed RequestState = 2
	Complete  RequestState = 3
)*/
