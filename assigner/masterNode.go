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

type StateMessage struct {
	State    e.ElevatorState
	SenderID int
}

type OrderMessage struct {
	Order      driver.ButtonEvent
	ReceiverID int
}

func MasterNode(
	ch_peerUpdate chan p.PeerUpdate, //maybe remove??
	ch_peerTxEnable chan bool, //maybe remove??
	ch_calculateNewOrder <-chan driver.ButtonEvent,
	ch_localStateToMaster <-chan StateMessage,
	ch_orderCalculated chan<- OrderMessage,
	ch_updateGlobalStateDemand chan<- []e.ElevatorState) {

	masterID := 0
	globalState := u.InitGlobalState()
	//println("LEGNTH GLOBALSTATE:", len(globalState))

	for {
		if masterID == c.ElevatorID {
			
			select {
			case order := <-ch_calculateNewOrder:
				fmt.Println("NEW ORDER TO MASTER:", order)
				lowestCostElevator := calculateCost(globalState, order)
				ch_orderCalculated <- OrderMessage{Order: order, ReceiverID: lowestCostElevator}

			case stateMsg:= <-ch_localStateToMaster:
				elevatorID := stateMsg.SenderID
				globalState[elevatorID] = stateMsg.State
				ch_updateGlobalStateDemand <- globalState
			}
		}
	}

}


func calculateCost(GlobalState []e.ElevatorState, order driver.ButtonEvent) int {
	var lowestCostID int
	cost := 9999

	for index, localState := range GlobalState {

		ElevatorCost := Cost(localState, order)
		//println("ID:", index, ", COST:", ElevatorCost)
		if ElevatorCost < cost {

			cost = ElevatorCost
			lowestCostID = index
		}
	}

	return lowestCostID
}
