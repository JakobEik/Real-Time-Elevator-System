package assigner

import (
	c "Project/config"
	"Project/driver"
	e "Project/elevator"
	p "Project/network/peers"
	u "Project/utils"
	"fmt"
	"math"

	//drv "Project/driver"
)

func Master(
	ch_peerUpdate chan p.PeerUpdate, //maybe remove??
	ch_peerTxEnable chan bool, //maybe remove??
	ch_messageToNetwork chan<- c.NetworkMessage,
	ch_messageFromNetwork <-chan c.NetworkMessage) {

	globalState := u.InitGlobalState()


	for {
		select {
		
		case message := <-ch_messageFromNetwork:
			switch message.MsgType {
			case c.NewOrder:
				newOrderEvent(message, globalState, ch_messageToNetwork)

			case c.OrderDone:
				orderDoneEvent(message, &globalState, ch_messageToNetwork)

			case c.MsgReceived:
				ElevatorID := message.SenderID

				UpdateGlobalState(&globalState, ElevatorID, message)
				msg := makeMessage(c.ToEveryone, globalState, c.ChangeYourState)
				ch_messageToNetwork <- msg

			case c.LocalStateChange:
				ElevatorID := message.SenderID
				UpdateGlobalState(&globalState, ElevatorID, message)
				fmt.Println(globalState[ElevatorID].Floor)
				//msg := makeMessage(c.ToEveryone, globalState, c.ChangeYourState)
				//ch_messageToNetwork <- msg
			}

		}

	}

}

func newOrderEvent(message c.NetworkMessage, globalState []e.ElevatorState, ch_messageToNetwork chan<- c.NetworkMessage){
	order := message.Msg.(driver.ButtonEvent)
	lowestCostElevator := calculateCost(globalState, order)
	msg := makeMessage(lowestCostElevator, order, c.DoOrder)
	ch_messageToNetwork <- msg
}

func orderDoneEvent(message c.NetworkMessage, globalState *[]e.ElevatorState, ch_messageToNetwork chan<- c.NetworkMessage){
	ElevatorID := message.SenderID
	UpdateGlobalState(globalState, ElevatorID, message)
	msg := makeMessage(c.ToEveryone, globalState, c.ChangeYourState)
	ch_messageToNetwork <- msg
}

func calculateCost(GlobalState []e.ElevatorState, order driver.ButtonEvent) int {

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

func UpdateGlobalState(globalState *[]e.ElevatorState, elevatorID int, message c.NetworkMessage) {
	state := message.Msg.(map[string]interface{})
	fmt.Println(state)
	var newState e.ElevatorState
	u.ConvertMapToStruct(state, &newState)
	
	(*globalState)[elevatorID] = newState
	
	/*
	if message.MsgType == c.OrderDone || message.MsgType == c.MsgReceived {
		order := message.Msg.(driver.ButtonEvent)
		if (*GlobalState)[elevatorID].Orders[order.Floor][order.Button] {
			(*GlobalState)[elevatorID].Orders[order.Floor][order.Button] = false
		} else {
			(*GlobalState)[elevatorID].Orders[order.Floor][order.Button] = true
		}
	} else if message.MsgType == c.LocalStateChange {
		state := message.Msg.(e.ElevatorState)
		(*GlobalState)[elevatorID] = state
	}*/
}
