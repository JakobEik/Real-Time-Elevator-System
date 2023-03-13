package utilities

import e "Project/elevator"

func InitGlobalState() []e.ElevatorState {
	globalState := make([]e.ElevatorState, N_ELEVATORS)
	for i := 0; i < N_ELEVATORS-1; i++ {
		globalState = append(globalState, e.InitElev())
	}
	return globalState
}

func CreateMessage(
	SenderID int,
	ReceiverID int,
	Msg any,
	MsgType MessageType) NetworkMessage {
	return NetworkMessage{
		SenderID:   SenderID,
		ReceiverID: ReceiverID,
		Msg:        Msg,
		MsgType:    MsgType}
}
