package assigner

import (
	c "Project/config"
	p "Project/network/peers"
	drv "Project/driver"
)

func master(ch_peerUpdate chan p.PeerUpdate, ch_peerTxEnable chan bool,
	//Tx Channels
	ch_TxGlobalState chan c.GlobalState, ch_TxNewOrder chan drv.ButtonEvent, ch_TxOrderDone chan drv.ButtonEvent,
	ch_TxOrderAccepted chan drv.ButtonEvent, ch_TxRequestGlobalState chan bool, ch_TxDoOrder chan drv.ButtonEvent,
	ch_TxChangeYourState chan TODO, ch_TxMsgReceived chan bool, 
	//Rx Channels
	
	 ){

}


func calculateCost(){



}


