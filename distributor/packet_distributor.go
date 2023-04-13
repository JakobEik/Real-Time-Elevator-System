package distributor

import (
	c "Project/config"
	"Project/utils"
	"time"
)

// Packet Fields need to be exported so network module can use them
type Packet struct {
	Msg            c.NetworkMessage
	SequenceNumber uint32
}

// Used for resending packets that havent been acknowledged
type packetWithAttempts struct {
	packet   Packet
	attempts uint8
}

// PacketDistributor distributes all messages from local to the network, and continue to send if no acknowledgment of
// the message has been received.
// It also distributes the messages from the network to the correct modules in the local elevator such that
// one module does not receive a message from the packet channel that is not intended for it, resulting in
// the message never arriving to the correct module.
func PacketDistributor(
	ch_packetFromNetwork <-chan Packet,
	ch_packetToNetwork chan<- Packet,
	ch_msgToPack <-chan c.NetworkMessage,
	ch_msgToAssigner, ch_msgToDistributor chan<- c.NetworkMessage) {

	var sequenceNumber uint32 = 0
	ticker := time.NewTicker(c.ConfirmationWaitDuration)
	packetsToAcknowledge := make(map[uint32]packetWithAttempts, 512) // map[seqNum]

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
					delete(packetsToAcknowledge, packet.SequenceNumber)
				}
			}

		case msg := <-ch_msgToPack:
			packet := pack(msg, sequenceNumber)
			packetsToAcknowledge[sequenceNumber] = packetWithAttempts{packet: packet, attempts: 1}
			ch_packetToNetwork <- packet
			sequenceNumber = incrementWithOverflow(sequenceNumber)
			// Function blocks until received or time limit reached. Might change this!
			//sendPacketUntilConfirmation(ch_packetToNetwork, ch_msgReceived, packet)

		case <-ticker.C:
			for seqNum, noneConfirmedPacket := range packetsToAcknowledge {
				if noneConfirmedPacket.attempts > c.NumOfRetries {
					delete(packetsToAcknowledge, seqNum)
				} else {
					noneConfirmedPacket.attempts++
					packetsToAcknowledge[seqNum] = noneConfirmedPacket
					ch_packetToNetwork <- noneConfirmedPacket.packet
				}
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
