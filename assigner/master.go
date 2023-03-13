package assigner

import (
	c "Project/config"
	p "Project/network/peers"
	//drv "Project/driver"
)

func master(
	ch_peerUpdate chan p.PeerUpdate, //maybe remove??
	ch_peerTxEnable chan bool,       //maybe remove??
	ch_messageToNetwork chan<- c.NetworkMessage,
	ch_messageFromNetwork <-chan c.NetworkMessage) {

	for {
		select {
		case message := <-ch_messageFromNetwork:
			switch message.MsgType {
			case c.GlobalState:

			case c.NewOrder:
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

func makeMessage(Receiver int, sender int, message any) c.NetworkMessage {

}
