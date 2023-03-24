package utils

import (
	"Project/config"
	e "Project/elevator"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func InitGlobalState() []e.ElevatorState {
	globalState := make([]e.ElevatorState, config.N_ELEVATORS)
	for i := 0; i < config.N_ELEVATORS; i++ {
		globalState[i] = e.InitElev(0)
	}
	return globalState
}

func CreatePacket(receiverID int, content any, msgType config.MessageType) config.Packet {
	msg := config.NetworkMessage{
		SenderID:   config.ElevatorID,
		MasterID:   0,
		ReceiverID: receiverID,
		Content:    content,
		MsgType:    msgType}
	
	return config.Packet{Msg: msg, Checksum: 0}	
}

func CastToType(data map[string]interface{}, myStruct interface{}) {

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
