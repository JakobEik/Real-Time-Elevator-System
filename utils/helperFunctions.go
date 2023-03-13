package utils

import (
	"Project/config"
	e "Project/elevator"
)

func InitGlobalState() []e.ElevatorState {
	globalState := make([]e.ElevatorState, config.N_ELEVATORS)
	for i := 0; i < config.N_ELEVATORS-1; i++ {
		globalState = append(globalState, e.InitElev())
	}
	return globalState
}

func CreateMessage(receiverID int, masterID int, msg any, msgType config.MessageType) config.NetworkMessage {
	return config.NetworkMessage{
		SenderID:   config.ElevatorID,
		MasterID:   masterID,
		ReceiverID: receiverID,
		Msg:        msg,
		MsgType:    msgType}
}
