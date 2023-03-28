package utils

import (
	"Project/config"
	e "Project/elevator"
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"log"

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

func CastToType(data interface{}, myStruct interface{}) {

	// Use mapstructure to map the data from the map to the struct
	conf := &mapstructure.DecoderConfig{
		ErrorUnused: true,
		Result:      &myStruct,
	}
	decoder, err := mapstructure.NewDecoder(conf)
	if err != nil {
		fmt.Println(err)
	}
	if err := decoder.Decode(data); err != nil {
		fmt.Println(err)
	}

}

func checksum(message any) uint32 {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(message)
	if err != nil {
		log.Fatal(err)
	}
	return crc32.ChecksumIEEE(buf.Bytes())
}
