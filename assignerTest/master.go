package assignerTest

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	p "Project/network/peers"
	"Project/utils"
	"fmt"
	"strconv"
)

type masterNode struct {
	ch_msgToPack chan<- c.NetworkMessage
	peersOnline  []int
	globalState  []e.ElevatorState
}

func (m masterNode) newOrderEvent(msg c.NetworkMessage) {
	if isElevatorOnline(msg.SenderID, m.peersOnline) {
		var order drv.ButtonEvent
		utils.DecodeContentToStruct(msg.Content, &order)
		bestScoreElevator := getBestElevatorForOrder(m.globalState, order, m.peersOnline)
		packet := utils.CreateMessage(bestScoreElevator, order, c.DO_ORDER)
		m.ch_msgToPack <- packet
		//fmt.Println("ORDER to elevator:", bestScoreElevator)
	}

}

func (m masterNode) newLocalStateEvent(msg c.NetworkMessage) {
	var state e.ElevatorState
	utils.DecodeContentToStruct(msg.Content, &state)
	elevatorID := msg.SenderID
	m.globalState[elevatorID] = state

	for _, ID := range m.peersOnline {
		globalStateUpdate := utils.CreateMessage(ID, m.globalState, c.UPDATE_GLOBAL_STATE)
		m.ch_msgToPack <- globalStateUpdate

		globalHallOrders := getGlobalHallOrders(m.globalState, m.peersOnline)
		hallLightsUpdate := utils.CreateMessage(ID, globalHallOrders, c.HALL_LIGHTS_UPDATE)
		m.ch_msgToPack <- hallLightsUpdate
	}

}

func (m masterNode) peerUpdate(update p.PeerUpdate) {
	println("THIS IS MASTER")
	m.peersOnline = stringArrayToIntArray(update.Peers)
	// Distribute orders from lost peers if there are any
	if len(update.Lost) > 0 {
		IDs := stringArrayToIntArray(update.Lost)
		for _, elevatorID := range IDs {
			distributeOrders(m.globalState[elevatorID], m.ch_msgToPack)
		}
	}
	// Update new peer
	if len(update.New) > 0 {
		elevatorID, _ := strconv.Atoi(update.New)
		stateUpdate := utils.CreateMessage(elevatorID, m.globalState, c.UPDATE_GLOBAL_STATE)
		m.ch_msgToPack <- stateUpdate

		hallOrders := getGlobalHallOrders(m.globalState, m.peersOnline)
		hallOrdersUpdate := utils.CreateMessage(elevatorID, hallOrders, c.HALL_LIGHTS_UPDATE)
		m.ch_msgToPack <- hallOrdersUpdate

		// Send the cab calls this peer had before connection loss
		cabCalls := getCabCalls(m.globalState[elevatorID])
		for _, order := range cabCalls {
			cabCallsUpdate := utils.CreateMessage(elevatorID, order, c.DO_ORDER)
			m.ch_msgToPack <- cabCallsUpdate
		}

	}

}

func isElevatorOnlineStr(ID int, peersOnlineStr []string) bool {
	peersOnline := stringArrayToIntArray(peersOnlineStr)
	return isElevatorOnline(ID, peersOnline)
}

func isElevatorOnline(ID int, peersOnline []int) bool {
	for _, peer := range peersOnline {
		if peer == ID {
			return true
		}
	}
	return false
}

func getCabCalls(e e.ElevatorState) []drv.ButtonEvent {
	orders := e.Orders
	var cabColumn []bool
	var cabCalls []drv.ButtonEvent
	for _, row := range orders {
		cabColumn = append(cabColumn, row[len(row)-1])
	}
	for floor, cabCall := range cabColumn {
		if cabCall {
			order := drv.ButtonEvent{Floor: floor, Button: drv.BT_Cab}
			cabCalls = append(cabCalls, order)
		}
	}
	return cabCalls
}

// distributeOrders sends every order from the given elevator to the masterNode of the system
func distributeOrders(elevator e.ElevatorState, ch_msgToPack chan<- c.NetworkMessage) {
	orders := elevator.Orders
	//println("DISTRIBUTE THIS ELEVATOR")
	//e.PrintState(elevator)
	for floor := range orders {
		for btn := 0; btn < c.N_BUTTONS-1; btn++ {
			if orders[floor][btn] == true {
				order := drv.ButtonEvent{Floor: floor, Button: drv.ButtonType(btn)}
				msg := utils.CreateMessage(c.MasterID, order, c.NEW_ORDER)
				//fmt.Println("DISTRIBUTE ORDER:", order)
				ch_msgToPack <- msg
			}
		}
	}
}

func getGlobalHallOrders(globalState []e.ElevatorState, onlineElevs []int) [][]bool {
	buttons := e.InitElev(0).Orders
	for _, ID := range onlineElevs {
		for floor := 0; floor < c.N_FLOORS; floor++ {
			for btn := 0; btn < c.N_BUTTONS-1; btn++ {
				if globalState[ID].Orders[floor][btn] == true {
					buttons[floor][btn] = true
				}
			}
		}
	}
	return buttons
}

func stringArrayToIntArray(strings []string) []int {
	ints := make([]int, len(strings))
	var err error
	for i, s := range strings {
		ints[i], err = strconv.Atoi(s)
		if err != nil {
			fmt.Println("Error:", err)
			panic(err)
		}
	}
	return ints
}
