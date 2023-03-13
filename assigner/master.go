package assigner

import (
	p "Project/network/peers"
	util "Project/utilities"

	//drv "Project/driver"
)

func master(
	ch_peerUpdate chan p.PeerUpdate, //maybe remove??
	ch_peerTxEnable chan bool,       //maybe remove??
	ch_messageToNetwork <-chan util.NetworkMessage,
	ch_messageFromNetwork chan<- util.NetworkMessage) {

	for {
		select {
		case message := <-ch_messageToNetwork:
			switch message.MsgType {
			case util.GlobalState:

			case util.NewOrder:
			case OrderDone:
			case OrderAccepted:
			case RequestGlobalState:
			case DoOrder:
			case ChangeYourState:
			case MsgReceived:

			}

		}

	}

}

func calculateCost() {

}

func makeMessage(Receiver int, sender int, message any) util.NetworkMessage {

}
