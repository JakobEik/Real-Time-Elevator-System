package assigner

import (
	c "Project/config"
	drv "Project/driver"
	"math"
)

const DOOR_OPEN_TIME = 3
const TRAVEL_TIME = 3

func getBestElevatorForOrder(GlobalState []c.ElevatorState, order drv.ButtonEvent, peersOnline []int) int {
	var bestElevatorID int
	bestScore := 9999999
	for _, elevID := range peersOnline {
		ElevatorScore := score(GlobalState[elevID], order)
		//println("ID:", elevID, ", SCORE:", ElevatorScore)
		if ElevatorScore < bestScore {
			bestScore = ElevatorScore
			bestElevatorID = elevID
		}
	}

	return bestElevatorID
}

func score(elev c.ElevatorState, newOrder drv.ButtonEvent) int {
	estimatedTime := 0
	// Loop through all floors and orders of the current elevator
	for floor := 0; floor < c.N_FLOORS; floor++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if elev.Orders[floor][btn] {
				// Calculate time to reach this floor from current floor
				distance := int(math.Abs(float64(elev.Floor - floor)))
				timeToReachFloor := distance * TRAVEL_TIME
				timeToServeOrder := DOOR_OPEN_TIME
				estimatedTime += timeToReachFloor + timeToServeOrder
			}
		}
	}
	distance := int(math.Abs(float64(elev.Floor - newOrder.Floor)))
	timeToReachFloor := distance * TRAVEL_TIME
	timeToServeOrder := DOOR_OPEN_TIME
	estimatedTime += timeToReachFloor + timeToServeOrder
	return estimatedTime
}
