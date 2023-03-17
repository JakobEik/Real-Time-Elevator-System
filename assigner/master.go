package assigner

import (
	c "Project/config"
	"Project/driver"
	"Project/elevator"
	p "Project/network/peers"
	u "Project/utils"
	"math"
	//drv "Project/driver"
)

func master(
	ch_peerUpdate chan p.PeerUpdate, //maybe remove??
	ch_peerTxEnable chan bool, //maybe remove??
	ch_messageToNetwork chan<- c.NetworkMessage,
	ch_messageFromNetwork <-chan c.NetworkMessage) {

	GlobalState := u.InitGlobalState()
	Gs_pt := &GlobalState

	for {
		select {
		case message := <-ch_messageFromNetwork:
			switch message.MsgType {
			case c.NewOrder:
				order := message.Msg.(driver.ButtonEvent)
				lowestCostElevator := calculateCost(GlobalState, order)
				msg := makeMessage(lowestCostElevator, order, c.DoOrder)
				ch_messageToNetwork <- msg

			case c.OrderDone:
				ElevatorID := message.SenderID

				UpdateGlobalState(Gs_pt, ElevatorID, message)
				msg := makeMessage(c.ToEveryone, GlobalState, c.ChangeYourState)
				ch_messageToNetwork <- msg

			case c.MsgReceived:
				ElevatorID := message.SenderID

				UpdateGlobalState(Gs_pt, ElevatorID, message)
				msg := makeMessage(c.ToEveryone, GlobalState, c.ChangeYourState)
				ch_messageToNetwork <- msg

			case c.LocalStateChange:
				ElevatorID := message.SenderID
				UpdateGlobalState(Gs_pt, ElevatorID, message)

				msg := makeMessage(c.ToEveryone, GlobalState, c.ChangeYourState)
				ch_messageToNetwork <- msg
			}

		}

	}

}

func calculateCost(GlobalState []elevator.ElevatorState, order driver.ButtonEvent) int {

	var lowestCostID int
	cost := int(math.Inf(1))

	for index, localState := range GlobalState {

		ElevatorCost := Cost(localState, order)
		if cost == int(math.Inf(1)) || ElevatorCost < cost {
			cost = ElevatorCost
			lowestCostID = index
		}
	}

	return lowestCostID
}

func makeMessage(receverID int, message any, messagetype c.MessageType) c.NetworkMessage {
	msg := c.NetworkMessage{}
	msg.SenderID = c.ElevatorID
	msg.MasterID = c.ElevatorID
	msg.ReceiverID = receverID
	msg.Msg = message
	msg.MsgType = messagetype

	return msg
}

func UpdateGlobalState(GlobalState *[]elevator.ElevatorState, elevatorID int, message c.NetworkMessage) {

	if message.MsgType == c.OrderDone || message.MsgType == c.MsgReceived {
		order := message.Msg.(driver.ButtonEvent)
		if (*GlobalState)[elevatorID].Orders[order.Floor][order.Button] {
			(*GlobalState)[elevatorID].Orders[order.Floor][order.Button] = false
		} else {
			(*GlobalState)[elevatorID].Orders[order.Floor][order.Button] = true
		}
	} else if message.MsgType == c.LocalStateChange {
		state := message.Msg.(elevator.ElevatorState)
		(*GlobalState)[elevatorID] = state
	}
}
