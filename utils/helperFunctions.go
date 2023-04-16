package utils

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func InitGlobalState() []c.ElevatorState {
	globalState := make([]c.ElevatorState, c.N_ELEVATORS)
	for i := 0; i < c.N_ELEVATORS; i++ {
		globalState[i] = e.InitElev(0)
	}
	return globalState
}

func CreateMessage(receiverID int, content any, msgType c.MessageType) c.NetworkMessage {
	msg := c.NetworkMessage{
		SenderID:   c.ElevatorID,
		ReceiverID: receiverID,
		Content:    content,
		Type:       msgType}

	return msg
}

// DecodeContent converts the data into the same struct it was sent as.
// Is needed since the content of the network messages are received as map[string]interface{}
func DecodeContent(oldMsg c.NetworkMessage) c.NetworkMessage {
	msg := oldMsg
	switch msg.Type {
	case c.DO_ORDER, c.NEW_ORDER:
		var order drv.ButtonEvent
		DecodeContentToStruct(msg.Content, &order)
		msg.Content = order

	case c.LOCAL_STATE_CHANGED:
		var state c.ElevatorState
		DecodeContentToStruct(msg.Content, &state)
		msg.Content = state

	case c.HALL_LIGHTS_UPDATE:
		var orders [][]bool
		DecodeContentToStruct(msg.Content, &orders)
		msg.Content = orders

	case c.NEW_MASTER, c.UPDATE_GLOBAL_STATE:
		globalState := make([]c.ElevatorState, c.N_ELEVATORS)
		content := msg.Content.([]interface{})
		// Iterates through the array, converts each one to ElevatorState and updates the global state
		for i, value := range content {
			var state c.ElevatorState
			DecodeContentToStruct(value, &state)
			globalState[i] = state
		}
		msg.Content = globalState
	default:
		panic("MESSAGE TYPE NOT IMPLEMENTED : " + msg.Type.String())
	}
	return msg
}

// DecodeContentToStruct uses the mapstructure package to convert one arbitrary Go type into another.
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
