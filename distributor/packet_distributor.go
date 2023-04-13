package distributor

import (
	c "Project/config"
	"Project/utils"
	"fmt"
	"time"
)

// Packet Fields need to be exported so network module can use them
type Packet struct {
	Msg            c.NetworkMessage
	SequenceNumber uint32
}

// Used for resending packets that havent been acknowledged
type packetWithResendAttempts struct {
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
	resendPacketsTicker := time.NewTicker(c.ConfirmationWaitDuration)
	packetsToBeConfirmed := make(map[uint32]packetWithResendAttempts, 512) // map[seqNum]

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
					delete(packetsToBeConfirmed, packet.SequenceNumber)
				}
			}

		case msg := <-ch_msgToPack:
			packet := pack(msg, sequenceNumber)
			packetsToBeConfirmed[sequenceNumber] = packetWithResendAttempts{packet: packet, attempts: 1}
			ch_packetToNetwork <- packet
			sequenceNumber = incrementWithOverflow(sequenceNumber)

		case <-resendPacketsTicker.C:
			// Resends all packets that have not been confirmed yet every x milliseconds
			for seqNum, packetAndAttempts := range packetsToBeConfirmed {
				print("RESENDING PACKETS, AMOUNT:", len(packetsToBeConfirmed))
				fmt.Println(", TYPE:", packetsToBeConfirmed[seqNum].packet.Msg.Type)
				if packetAndAttempts.attempts > c.NumOfRetries {
					delete(packetsToBeConfirmed, seqNum)
				} else {
					packetAndAttempts.attempts++
					packetsToBeConfirmed[seqNum] = packetAndAttempts
					ch_packetToNetwork <- packetAndAttempts.packet
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
