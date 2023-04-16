package distributor

import (
	c "Project/config"
	"Project/utils"
	"Project/watchdog"
	"time"
)

// Packet the packet sendt over the network
type Packet struct {
	Msg            c.NetworkMessage
	SequenceNumber uint32
}

// Used for resending packets that havent been acknowledged yet
type packetWithResendAttempts struct {
	packet   Packet
	attempts uint8
}

const numOfRetries = 10
const confirmationWaitDuration = time.Millisecond * 50

// PacketDistributor distributes all packets between the network and the local modules
// If no acknowledgment for a packet has been received, this module will continue to send it until has been received.
// There is one channel each for assigner and distributor, this is to prevent a data race
func PacketDistributor(
	ch_packetFromNetwork <-chan Packet,
	ch_packetToNetwork chan<- Packet,
	ch_packMessage <-chan c.NetworkMessage,
	ch_msgToAssigner, ch_msgToDistributor chan<- c.NetworkMessage) {

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Packet Distributor")

	var sequenceNumber uint32 = 0
	resendPacketsTicker := time.NewTicker(confirmationWaitDuration)
	packetsToBeConfirmed := make(map[uint32]packetWithResendAttempts, 512) // map[seqNum]

	for {
		select {
		//Watchdog
		case <-ch_bark:
			ch_pet <- true
		//_________________FROM NETWORK________________
		case packet := <-ch_packetFromNetwork:
			if acceptPacket(packet) {
				switch packet.Msg.Type {
				// TO ASSIGNER
				case c.NEW_MASTER, c.LOCAL_STATE_CHANGED, c.UPDATE_GLOBAL_STATE, c.NEW_ORDER:
					ch_msgToAssigner <- packet.Msg
					ch_packetToNetwork <- confirmMessage(packet)

				// TO DISTRIBUTOR
				case c.DO_ORDER, c.HALL_LIGHTS_UPDATE:
					ch_msgToDistributor <- packet.Msg
					ch_packetToNetwork <- confirmMessage(packet)

				case c.MSG_RECEIVED:
					delete(packetsToBeConfirmed, packet.SequenceNumber)
				}
			}
		//_________________TO NETWORK________________
		case msg := <-ch_packMessage:
			packet := packMessage(msg, sequenceNumber)
			packetsToBeConfirmed[sequenceNumber] = packetWithResendAttempts{packet: packet, attempts: 1}
			ch_packetToNetwork <- packet
			sequenceNumber = incrementWithOverflow(sequenceNumber)

		case <-resendPacketsTicker.C:
			// Resends all packets that have not been confirmed yet every x milliseconds
			for seqNum, packetAndAttempts := range packetsToBeConfirmed {
				if packetAndAttempts.attempts > numOfRetries {
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
	return receiverID == c.ElevatorID

}

func packMessage(message c.NetworkMessage, seqNum uint32) Packet {
	packet := Packet{Msg: message, SequenceNumber: seqNum}
	return packet
}

func confirmMessage(packetReceived Packet) Packet {
	receiverID := packetReceived.Msg.SenderID
	msg := utils.CreateMessage(receiverID, true, c.MSG_RECEIVED)
	return packMessage(msg, packetReceived.SequenceNumber)
}
