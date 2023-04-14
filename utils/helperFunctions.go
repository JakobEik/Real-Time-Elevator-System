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

func CreateMessage(receiverID int, content any, msgType config.MessageType) config.NetworkMessage {
	msg := config.NetworkMessage{
		SenderID:   config.ElevatorID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType}

	return msg
}

// DecodeContentToStruct uses the mapstructure package to convert one arbitrary Go type into another.
// Is needed since the content of the network messages are received as map[string]interface{}, and
// this function converts the data into the same struct it was sent as.
func DecodeContentToStruct(data interface{}, correctStruct interface{}) {
	// Use mapstructure to map the data from the map to the struct
	conf := &mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &correctStruct,
	}
	decoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		fmt.Println("DECODE TO STRUCT ERROR:", err)
	}
	if err := decoder.Decode(data); err != nil {
		fmt.Println("DECODE TO STRUCT ERROR:", err)
	}

}

