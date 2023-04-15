package assigner2

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"Project/network/peers"
	"Project/utils"
	"fmt"
	"strconv"
)

func isElevatorOffline(ID int, peersOnline []int) bool {
	for _, peer := range peersOnline {
		if peer == ID {
			return false
		}
	}
	return true
}

// Checks if the elevator is new on the network or not
func isNewElevator(ID int, update peers.PeerUpdate) bool {
	peersOnline := stringArrayToIntArray(update.Peers)
	if len(update.New) > 0 {
		newElevID, _ := strconv.Atoi(update.New)
		if ID == newElevID {
			return false
		}
	}
	return !isElevatorOffline(ID, peersOnline)
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

// getMaster returns the elevator with the lowest ID
func getMaster(peersOnline []int) int {
	if len(peersOnline) == 0 {
		panic("NO MORE ELEVATORS ON NETWORK")
	}
	masterID := peersOnline[0]
	for _, elev := range peersOnline[1:] {
		if elev < masterID {
			masterID = elev
		}
	}
	return masterID
}
