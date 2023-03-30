package distributor

import (
	c "Project/config"
	"Project/utils"
	"time"
)

// Fields need to be exported so network module can use them
type Packet struct {
	Msg            c.NetworkMessage
	SequenceNumber uint32
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

	var sequenceNumber uint32 = 0
	ch_msgReceived := make(chan uint32, 512)

	for {
		select {
		case packet := <-ch_packetFromNetwork:
			//fmt.Println("RECEIVE:", packet.Msg.Type)
			if acceptPacket(packet) {
				switch packet.Msg.Type {
				// TO ASSIGNER
				case c.LOCAL_STATE_CHANGED:
					fallthrough
				case c.UPDATE_GLOBAL_STATE:
					fallthrough
				case c.NEW_ORDER:
					ch_msgToAssigner <- packet.Msg
					ch_packetToNetwork <- confirmMessage(packet)

				// TO DISTRIBUTOR
				case c.DO_ORDER:
					fallthrough
				case c.GLOBAL_HALL_ORDERS:
					ch_msgToDistributor <- packet.Msg
					ch_packetToNetwork <- confirmMessage(packet)

				case c.MSG_RECEIVED:
					ch_msgReceived <- packet.SequenceNumber

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
				//fmt.Println("SEND:", packet.Msg.Type)
				count++
			} else {
				//fmt.Println("PACKET FAILED TO RECEIVE CONFIRMATION")
				return
			}

		case seqNum := <-ch_msgReceived:
			if seqNum == packet.SequenceNumber {
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
	receiverID := packet.Msg.ReceiverID
	return receiverID == c.ElevatorID || receiverID == c.ToEveryone

}

func pack(message c.NetworkMessage, seqNum uint32) Packet {
	packet := Packet{Msg: message, SequenceNumber: seqNum}
	return packet
}

func confirmMessage(packetReceived Packet) Packet {
	receiverID := packetReceived.Msg.SenderID
	msg := utils.CreateMessage(receiverID, true, c.MSG_RECEIVED)
	return pack(msg, packetReceived.SequenceNumber)
}
