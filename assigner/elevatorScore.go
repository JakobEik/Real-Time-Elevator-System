package assigner

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"math"
)

const DOOR_OPEN_TIME = 3
const TRAVEL_TIME = 3

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

func chooseElevDirection(e e.ElevatorState) (drv.MotorDirection, c.Behavior) {
	switch e.Direction {
	case drv.MD_Up:
		if ordersAbove(e) {
			return drv.MD_Up, c.MOVING
		} else if ordersHere(e) {
			return drv.MD_Stop, c.DOOR_OPEN
		} else if ordersBelow(e) {
			return drv.MD_Down, c.MOVING
		} else {
			return drv.MD_Stop, c.IDLE
		}
	case drv.MD_Down:
		if ordersBelow(e) {
			return drv.MD_Down, c.MOVING
		} else if ordersHere(e) {
			return drv.MD_Up, c.DOOR_OPEN
		} else if ordersAbove(e) {
			return drv.MD_Up, c.MOVING
		} else {
			return drv.MD_Stop, c.IDLE
		}
	case drv.MD_Stop:
		if ordersHere(e) {
			return drv.MD_Stop, c.DOOR_OPEN
		} else if ordersAbove(e) {
			return drv.MD_Up, c.MOVING
		} else if ordersBelow(e) {
			return drv.MD_Down, c.MOVING
		} else {
			return drv.MD_Stop, c.IDLE
		}
	default:
		return drv.MD_Stop, c.IDLE
	}
}

func score1(e_old e.ElevatorState, b drv.ButtonEvent) int {
	e := e_old
	floor := b.Floor
	button := b.Button
	e.Orders[floor][button] = true

	var arrivedAtRequest bool

	ifEqual := func(inner_b drv.ButtonType, inner_f int) {
		if inner_b == b.Button {
			arrivedAtRequest = true
		}
	}

	duration := 0

	switch e.Behavior {
	case c.IDLE:
		e.Direction, e.Behavior = chooseElevDirection(e)
		if e.Direction == drv.MD_Stop {
			return duration
		}
	case c.MOVING:
		duration += TRAVEL_TIME / 2
		e.Floor += int(e.Direction)
	case c.DOOR_OPEN:
		duration -= DOOR_OPEN_TIME / 2
	}

	for {
		if shouldStop(e) {
			e = requests_clearAtCurrentFloor(e, ifEqual)
			if arrivedAtRequest {
				return duration
			}
			duration += DOOR_OPEN_TIME
			e.Direction, e.Behavior = chooseElevDirection(e)
		}
		e.Floor += int(e.Direction)
		duration += TRAVEL_TIME
	}
}

func requests_clearAtCurrentFloor(e_old e.ElevatorState, onClearedRequest func(btn drv.ButtonType, floor int)) e.ElevatorState {
	e := e_old
	// This shouldn't clear every single order - just to make the example shorter
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.Orders[e.Floor][btn] {
			e.Orders[e.Floor][btn] = false
			if onClearedRequest != nil {
				onClearedRequest(drv.ButtonType(btn), e.Floor)
			}
		}
	}
	return e
}

func ordersAbove(e e.ElevatorState) bool {
	if e.Floor >= c.N_FLOORS {
		return false
	}
	for f := e.Floor + 1; f < c.N_FLOORS; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.Orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(e e.ElevatorState) bool {
	if e.Floor <= 0 {
		return false
	}
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.Orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersHere(e e.ElevatorState) bool {
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.Orders[e.Floor][btn] {
			return true
		}
	}
	return false
}

func shouldStop(e e.ElevatorState) bool {
	switch e.Direction {
	case drv.MD_Down:
		return e.Orders[e.Floor][drv.BT_HallDown] ||
			e.Orders[e.Floor][drv.BT_Cab] ||
			!ordersBelow(e)
	case drv.MD_Up:

		return e.Orders[e.Floor][drv.BT_HallUp] ||
			e.Orders[e.Floor][drv.BT_Cab] ||
			!ordersAbove(e)
	case drv.MD_Stop:
		fallthrough
	default:
		return true
	}

}
