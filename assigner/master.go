package assigner

import (
	c "Project/config"
	"Project/driver"
	e "Project/elevator"
	p "Project/network/peers"
	u "Project/utils"
	"fmt"
	//drv "Project/driver"
)

func Master(
	ch_peerUpdate chan p.PeerUpdate, //maybe remove??
	ch_peerTxEnable chan bool, //maybe remove??
	ch_messageToNetwork chan<- c.NetworkMessage,
	ch_messageFromNetwork <-chan c.NetworkMessage) {

	masterID := 0
	globalState := u.InitGlobalState()
	//println("LEGNTH GLOBALSTATE:", len(globalState))

	for {
		select {
		case message := <-ch_messageFromNetwork:
			if masterID == c.ElevatorID {
				switch message.MsgType {
				case c.NewOrder:
					fmt.Println("NEW ORDER:", message.Content)
					newOrderEvent(message, globalState, ch_messageToNetwork)

				case c.OrderDone:
					orderDoneEvent(message, &globalState, ch_messageToNetwork)

				case c.MsgReceived:
					ElevatorID := message.SenderID

					UpdateGlobalState(&globalState, ElevatorID, message)
					msg := u.CreateMessage(c.ToEveryone, c.ElevatorID, globalState, c.ChangeYourState)
					ch_messageToNetwork <- msg

				case c.LocalStateChange:
					ElevatorID := message.SenderID
					UpdateGlobalState(&globalState, ElevatorID, message)
					//msg := makeMessage(c.ToEveryone, globalState, c.ChangeYourState)
					//ch_messageToNetwork <- msg
				}
			}

		}

	}

}

func newOrderEvent(message c.NetworkMessage, globalState []e.ElevatorState, ch_messageToNetwork chan<- c.NetworkMessage) {
	content := message.Content.(map[string]interface{})
	var order driver.ButtonEvent
	u.ConvertMapToStruct(content, &order)
	lowestCostElevator := calculateCost(globalState, order, message.SenderID)
	println("LOWEST COST ELEVATOR:", lowestCostElevator)
	msg := u.CreateMessage(lowestCostElevator, c.ElevatorID, order, c.DoOrder)
	ch_messageToNetwork <- msg
}

func orderDoneEvent(message c.NetworkMessage, globalState *[]e.ElevatorState, ch_messageToNetwork chan<- c.NetworkMessage) {
	ElevatorID := message.SenderID
	UpdateGlobalState(globalState, ElevatorID, message)
	msg := u.CreateMessage(c.ToEveryone, c.ElevatorID, globalState, c.ChangeYourState)
	ch_messageToNetwork <- msg
}

func calculateCost(GlobalState []e.ElevatorState, order driver.ButtonEvent, senderID int) int {
	if order.Button == driver.BT_Cab {
		return senderID
	}
	var lowestCostID int
	cost := 9999

	for index, localState := range GlobalState {

		ElevatorCost := Cost(localState, order)
		println("ID:", index, ", COST:", ElevatorCost)
		if ElevatorCost < cost {

			cost = ElevatorCost
			lowestCostID = index
		}
	}

	return lowestCostID
}

func UpdateGlobalState(globalState *[]e.ElevatorState, elevatorID int, message c.NetworkMessage) {
	state := message.Content.(map[string]interface{})
	var newState e.ElevatorState
	u.ConvertMapToStruct(state, &newState)

	(*globalState)[elevatorID] = newState
	switch message.MsgType {
	case c.LocalStateChange:
		content := message.Content.(map[string]interface{})
		var state e.ElevatorState
		u.ConvertMapToStruct(content, &state)
		(*globalState)[elevatorID] = state

	}

	/*if message.MsgType == c.OrderDone || message.MsgType == c.MsgReceived {
		order := message.Content.(driver.ButtonEvent)
		if (*globalState)[elevatorID].Orders[order.Floor][order.Button] {
			(*globalState)[elevatorID].Orders[order.Floor][order.Button] = false
		} else {
			(*globalState)[elevatorID].Orders[order.Floor][order.Button] = true
		}
	} else if message.MsgType == c.LocalStateChange {
		state := message.Content.(e.ElevatorState)
		(*globalState)[elevatorID] = state
	}*/
}
