package assigner2

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"math"
)

const DOOR_OPEN_TIME = 3
const TRAVEL_TIME = 3

func getBestElevatorForOrder1(GlobalState []e.ElevatorState, order drv.ButtonEvent, peersOnline []int) int {
	var bestElevatorID int
	bestScore := -9999
	for _, elevID := range peersOnline {
		ElevatorScore := score(GlobalState[elevID], order)
		//println("ID:", elevID, ", SCORE:", ElevatorScore)
		if ElevatorScore > bestScore {
			bestScore = ElevatorScore
			bestElevatorID = elevID
		}
	}

	return bestElevatorID
}

// Nearest Car Algorithm
func score(elev e.ElevatorState, order drv.ButtonEvent) int {

	ordBtn := order.Button
	score := 0
	dir := elev.Direction
	ordFloor := order.Floor
	eFloor := elev.Floor
	distance := int(math.Abs(float64(eFloor - ordFloor)))

	switch elev.Behavior {
	case c.DOOR_OPEN:
		fallthrough
	case c.IDLE:
		score = c.N_FLOORS + 1 - distance // score = N + 1 - d	Nearest car algorithm
	case c.MOVING:
		if (dir == drv.MD_Up && ordFloor > eFloor && ordBtn == drv.BT_HallUp) ||
			(dir == drv.MD_Down && ordFloor < eFloor && ordBtn == drv.BT_HallDown) {

			score = c.N_FLOORS + 2 - distance
		} else if (dir == drv.MD_Up && ordFloor > eFloor && ordBtn == drv.BT_HallDown) ||
			(dir == drv.MD_Down && ordFloor < eFloor && ordBtn == drv.BT_HallUp) {

			score = c.N_FLOORS + 1 - distance
		} else {
			score = 1
		}
	}
	switch order.Button {
	case drv.BT_HallDown:
		if elev.Orders[ordFloor][drv.BT_HallUp] == true {
			score = 0
		}
	case drv.BT_HallUp:
		if elev.Orders[ordFloor][drv.BT_HallDown] == true {
			score = 0
		}
	}
	return score - orderCount(elev)/2
}

func orderCount(e e.ElevatorState) int {
	orders := e.Orders
	count := 0
	for _, row := range orders {
		for _, element := range row {
			if element == true {
				count++
			}
		}
	}
	return count
}

func getBestElevatorForOrder(GlobalState []e.ElevatorState, order drv.ButtonEvent, peersOnline []int) int {
	var bestElevatorID int
	bestScore := 9999
	for _, elevID := range peersOnline {
		ElevatorScore := score1(GlobalState[elevID], order)
		//println("ID:", elevID, ", SCORE:", ElevatorScore)
		if ElevatorScore < bestScore {
			bestScore = ElevatorScore
			bestElevatorID = elevID
		}
	}

	return bestElevatorID
}

func score1(elev e.ElevatorState, newOrder drv.ButtonEvent) int {
	estimatedTime := 0

	// Loop through all floors and orders of the current elevator
	for floor := 0; floor < c.N_FLOORS; floor++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if elev.Orders[floor][btn] {
				// Calculate time to reach this floor from current floor
				distance := int(math.Abs(float64(elev.Floor - floor)))
				timeToReachFloor := distance * TRAVEL_TIME // 3 seconds per floor

				// Add time to open
				timeToServeOrder := DOOR_OPEN_TIME // 3 seconds to open

				// Update estimated time to complete pending orders
				estimatedTime += timeToReachFloor + timeToServeOrder
			}
		}
	}

	// Calculate estimated time to reach new order floor from current floor
	distance := int(math.Abs(float64(elev.Floor - newOrder.Floor)))
	timeToReachFloor := distance * TRAVEL_TIME // 3 seconds per floor

	// Add time to open
	timeToServeOrder := DOOR_OPEN_TIME // 3 seconds to open

	// Add estimated time for the new order to the total estimated time
	estimatedTime += timeToReachFloor + timeToServeOrder

	return estimatedTime
}
