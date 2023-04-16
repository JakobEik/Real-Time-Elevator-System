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
	if len(update.New) > 0 {
		newElevID, _ := strconv.Atoi(update.New)
		if ID == newElevID {
			return true
		}
	}
	return false
}

func getCabCalls(e c.ElevatorState) []drv.ButtonEvent {
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

// updateGlobalState sets the hall orders to false for the lost peers in the global state
// after the orders are distributed
func updateGlobalState(globalState []c.ElevatorState, lostPeers []int) []c.ElevatorState {
	elevators := filterElevator(globalState, lostPeers)
	for _, elevator := range elevators {
		elevator.Orders = removeHallOrders(elevator.Orders)
	}
	for index, ID := range lostPeers {
		globalState[ID] = elevators[index]
	}
	return globalState

}

func makeMessagesFromOrders(orders [][]bool) []c.NetworkMessage {
	var messages []c.NetworkMessage
	for floor, floors := range orders {
		for btn, orderValue := range floors {
			if orderValue {
				order := drv.ButtonEvent{Floor: floor, Button: drv.ButtonType(btn)}
				msg := utils.CreateMessage(c.MasterID, order, c.NEW_ORDER)
				messages = append(messages, msg)
			}
		}
	}
	return messages
}

// getCombinedOrders Returns the combined orders for the input elevators
func getCombinedOrders(elevators []c.ElevatorState) [][]bool {
	if len(elevators) == 0 {
		return nil
	}
	orders := e.InitElev(0).Orders
	for _, elevator := range elevators {
		for floor, floors := range elevator.Orders {
			for btn, order := range floors {
				orders[floor][btn] = orders[floor][btn] || order
			}
		}
	}
	return orders
}

func filterElevator(globalState []c.ElevatorState, elevator_IDs []int) []c.ElevatorState {
	var result []c.ElevatorState
	for _, ID := range elevator_IDs {
		result = append(result, globalState[ID])
	}
	return result
}

func getHallOrders(globalState []c.ElevatorState, IDs []int) [][]bool {
	elevators := filterElevator(globalState, IDs)
	globalOrders := getCombinedOrders(elevators)
	// Sets all Cab orders to false since these are not used
	for i := range globalOrders {
		globalOrders[i][len(globalOrders[i])-1] = false
	}
	return globalOrders
}

func removeHallOrders(orders [][]bool) [][]bool {
	for floor := range orders {
		for btn := range orders[floor] {
			if btn != len(orders[floor])-1 {
				orders[floor][btn] = false
			}
		}
	}
	return orders

}

func stringArrayToIntArray(strings []string) []int {
	result := make([]int, len(strings))
	for i, str := range strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("EROOR IN stringArrayToIntArray")
			panic(err)
		}
		result[i] = num
	}
	return result
}

// getMaster returns the elevator with the lowest ID
func getMaster(peersOnline []int) int {
	if len(peersOnline) == 0 {
		return c.ElevatorID
	}
	masterID := peersOnline[0]
	for _, elev := range peersOnline[1:] {
		if elev < masterID {
			masterID = elev
		}
	}
	return masterID
}
