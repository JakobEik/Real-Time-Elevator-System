package utilities

import e "Project/elevator"

func InitGlobalState() []e.ElevatorState {
	globalState := make([]e.ElevatorState, N_ELEVATORS)
	for i := 0; i < N_ELEVATORS-1; i++ {
		globalState = append(globalState, e.InitElev())
	}
	return globalState
}

func CreateMessage(receiverID int, masterID int, msg any, msgType MessageType) NetworkMessage {
	return NetworkMessage{
		SenderID:   ElevatorID,
		MasterID:   masterID,
		ReceiverID: receiverID,
		Msg:        msg,
		MsgType:    msgType}
}
