package utils

import (
	"Project/config"
	e "Project/elevator"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func InitGlobalState() []e.ElevatorState {
	globalState := make([]e.ElevatorState, config.N_ELEVATORS)
	for i := 0; i < config.N_ELEVATORS-1; i++ {
		globalState = append(globalState, e.InitElev())
	}
	return globalState
}

func CreateMessage(receiverID int, masterID int, content any, msgType config.MessageType) config.NetworkMessage {
	return config.NetworkMessage{
		SenderID:   config.ElevatorID,
		MasterID:   masterID,
		ReceiverID: receiverID,
		Content:    content,
		MsgType:    msgType}
}

func ConvertMapToStruct(data map[string]interface{}, myStruct interface{}) {

	// Use mapstructure to map the data from the map to the struct
	config := &mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &myStruct,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		fmt.Println(err)
	}
	if err := decoder.Decode(data); err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("%+v\n", myStruct)
	//fmt.Printf("t1: %T\n", myStruct)

}
