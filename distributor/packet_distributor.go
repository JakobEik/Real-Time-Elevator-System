package distributor

import (
	c "Project/config"
	"Project/driver"
	"Project/elevator"
	"Project/utils"
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"log"
	"time"
)

type Packet struct {
	msg            c.NetworkMessage
	sequenceNumber uint32
	checksum       uint32
}

// PacketDistributor distributes all messages from local to the network, and continue to send if no confirmation has
// been received.
// It also distributes the messages from the network to the correct modules in the local elevator such that
// one module does not receive a message from the packet channel that is not intended for it, resulting in
// the message never arriving to the correct module.
// The module also handles checksum and sends confirmation back to the sender
func PacketDistributor(
	ch_packetFromNetwork <-chan Packet,
	ch_packetToNetwork chan<- Packet,
	ch_msgToPack <-chan c.NetworkMessage,
	ch_msgToAssigner, ch_msgToDistributor chan<- c.NetworkMessage) {

	initGob()

	var sequenceNumber uint32 = 0
	ch_msgReceived := make(chan uint32, 512)

	for {
		select {
		case packet := <-ch_packetFromNetwork:
			fmt.Println("RECEIVE:", packet)
			if acceptPacket(packet) {
				switch packet.msg.Type {
				// TO ASSIGNER
				case c.LOCAL_STATE_CHANGED:
					fallthrough
				case c.NEW_ORDER:
					if c.MasterID == c.ElevatorID {
						ch_msgToAssigner <- packet.msg
						ch_packetToNetwork <- confirmMessage(packet)
					}
				// TO DISTRIBUTOR
				case c.DO_ORDER:
					fallthrough
				case c.GLOBAL_HALL_ORDERS:
					fallthrough
				case c.UPDATE_GLOBAL_STATE:
					ch_msgToDistributor <- packet.msg
					ch_packetToNetwork <- confirmMessage(packet)

				case c.MSG_RECEIVED:
					ch_msgReceived <- packet.sequenceNumber

				}
			}

		case msg := <-ch_msgToPack:
			packet := pack(msg, sequenceNumber)
			sequenceNumber = incrementWithOverflow(sequenceNumber)
			// Function blocks until received or time limit reached. Might change this!
			sendPacketUntilConfirmation(ch_packetToNetwork, ch_msgReceived, packet)

		}
	}

}

func sendPacketUntilConfirmation(ch_packetToNetwork chan<- Packet, ch_msgReceived <-chan uint32, packet Packet) {
	ticker := time.NewTicker(c.ConfirmationWaitDuration)

	// stop the ticker when the function returns
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-ticker.C:
			if count < c.NumOfRetries {
				ch_packetToNetwork <- packet
				fmt.Println("SEND:", packet)
				count++
			} else {
				fmt.Println("PACKET FAILED TO RECEIVE CONFIRMATION")
				return
			}

		case seqNum := <-ch_msgReceived:
			if seqNum == packet.sequenceNumber {
				return
			}
		}
	}
}

// removeIfPresent removes a uint32 value from a map if it is present when it receives it on a channel
// MIGHT USE LATER
func removeIfPresent(c <-chan uint32, m map[uint32]bool) {
	for {
		select {
		case n := <-c: // when a value is received on the channel
			if m[n] { // if the map contains the value
				delete(m, n) // delete the value from the map
			}
		}
	}
}

func incrementWithOverflow(number uint32) uint32 {
	if number > 4000000000 {
		return 0
	} else {
		return number + 1
	}

}

func acceptPacket(packet Packet) bool {
	checksumCorrect := packet.checksum == checksum(packet.msg)
	//if !checksumCorrect {
	//	println("WRONG CHECKSUM! Calculated:", checksum(packet.msg), ", received:", packet.checksum)
	//}
	receiverID := packet.msg.ReceiverID
	return checksumCorrect && (receiverID == c.ElevatorID || receiverID == c.ToEveryone)

}

func pack(message c.NetworkMessage, seqNum uint32) Packet {
	return Packet{msg: message, sequenceNumber: seqNum, checksum: checksum(message)}
}

func confirmMessage(packetReceived Packet) Packet {
	receiverID := packetReceived.msg.SenderID
	msg := utils.CreateMessage(receiverID, true, c.MSG_RECEIVED)
	return pack(msg, packetReceived.sequenceNumber)
}

func checksum(message any) uint32 {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(message)
	if err != nil {
		log.Fatal("CHECKSUM ERROR:", err)
	}
	return crc32.ChecksumIEEE(buf.Bytes())
}

func initGob() {
	gob.Register(elevator.ElevatorState{})
	gob.Register(driver.ButtonEvent{})
}
