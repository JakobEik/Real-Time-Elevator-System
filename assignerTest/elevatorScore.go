package assignerTest

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"math"
)

func getBestElevatorForOrder(GlobalState []e.ElevatorState, order drv.ButtonEvent, elevatorIDs []int) int {
	var bestElevatorID int
	bestScore := -9999
	for _, elevID := range elevatorIDs {
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
