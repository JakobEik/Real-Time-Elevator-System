package assigner

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"math"
)

func Cost3(state e.ElevatorState, order drv.ButtonEvent) int {
	TRAVEL_TIME := 10
	elev := state

	if elev.Orders[order.Floor][order.Button] == true {
		return 0
	}
	elev.Orders[order.Floor][order.Button] = true

	duration := 0

	switch elev.Behavior {
	case c.Idle:
		dir, _ := chooseElevDirection(&elev)
		elev.Direction = dir
		if elev.Direction == drv.MD_Stop {
			return duration
		}

	case c.Moving:
		duration += TRAVEL_TIME / 2
		elev.Floor += int(elev.Direction)
	case c.DoorOpen:
		duration -= 2
	}

	for {
		if shouldStop(&elev) {
			clearAtCurrentFloor(&elev)

			if arrivedAtOrder(elev, order) {
				return duration
			}
			duration += 3
			dir, _ := chooseElevDirection(&elev)
			elev.Direction = dir
		}
		elev.Floor += int(elev.Direction)
		duration += TRAVEL_TIME
	}
}

func arrivedAtOrder(elev e.ElevatorState, order drv.ButtonEvent) bool {
	sameFloor := elev.Floor == order.Floor
	sameButton := (elev.Direction == drv.MD_Up && order.Button == drv.BT_HallUp) || (elev.Direction == drv.MD_Down && order.Button == drv.BT_HallDown)
	if order.Button == drv.BT_Cab {
		sameButton = true
	}
	return sameFloor && sameButton
}

// Nearest Car Algorithm
func Cost1(elev e.ElevatorState, order drv.ButtonEvent) int {

	cost := 0.0

	distance := elev.Floor - order.Floor

	if elev.Behavior == c.Idle && elev.Floor == order.Floor {
		return -10000
	}

	if distance/int(elev.Direction) >= 0 {
		if (elev.Direction == drv.MD_Up && order.Button == drv.BT_HallUp) ||
			(elev.Direction == drv.MD_Down && order.Button == drv.BT_HallDown) {
			cost = c.N_FLOORS + 2 - math.Abs(float64(distance))
		} else {
			cost = c.N_FLOORS + 1 - math.Abs(float64(distance))

		}

	} else {
		cost = 1
	}

	return -int(cost)
}

// Nearest Car Algorithm
func Cost(elev e.ElevatorState, order drv.ButtonEvent) int {

	ordBtn := order.Button
	cost := 0
	dir := elev.Direction
	ordFloor := order.Floor
	eFloor := elev.Floor
	distance := int(math.Abs(float64(eFloor - ordFloor)))

	switch elev.Behavior {
	case c.DoorOpen:
		fallthrough
	case c.Idle:
		cost = c.N_FLOORS + 1 - distance // cost = N + 1 - d	Nearest car algorithm
	case c.Moving:
		if (dir == drv.MD_Up && ordFloor > eFloor && ordBtn == drv.BT_HallUp) ||
			(dir == drv.MD_Down && ordFloor < eFloor && ordBtn == drv.BT_HallDown) {

			cost = c.N_FLOORS + 2 - distance
		} else if (dir == drv.MD_Up && ordFloor > eFloor && ordBtn == drv.BT_HallDown) ||
			(dir == drv.MD_Down && ordFloor < eFloor && ordBtn == drv.BT_HallUp) {

			cost = c.N_FLOORS + 1 - distance
		} else {
			cost = 1
		}
	}

	return -cost
}

/*
func Cost2(state e.ElevatorState, order drv.ButtonEvent) int { // Maybe make more efficient algorithm, the one in resources file??

		currFloor := state.Floor
		ordFloor := order.Floor
		ordBtn := order.Button
		var cost = 0

		distance := int(math.Abs(float64(currFloor) - float64(ordFloor)))

		if state.Behavior != c.Unavailable {
			switch state.Behavior {
			case c.Idle:
				cost = c.N_FLOORS + 1 - distance // cost = N + 1 - d	Nearest car algorithm
			case c.Moving:
				if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == drv.BT_HallUp) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == drv.BT_HallDown) {
					cost = c.N_FLOORS + 2 - distance
				} else if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == drv.BT_HallDown) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == drv.BT_HallUp) {
					cost = c.N_FLOORS + 1 - distance
				} else {
					cost = 1
				}
			}

		}

		return cost
	}
*/
func chooseElevDirection(e *e.ElevatorState) (drv.MotorDirection, c.Behavior) {
	switch e.Direction {
	case drv.MD_Up:
		if ordersAbove(e) {
			return drv.MD_Up, c.Moving
		} else if ordersHere(e) {
			return drv.MD_Stop, c.DoorOpen
		} else if ordersBelow(e) {
			return drv.MD_Down, c.Moving
		} else {
			return drv.MD_Stop, c.Idle
		}
	case drv.MD_Down:
		if ordersBelow(e) {
			return drv.MD_Down, c.Moving
		} else if ordersHere(e) {
			return drv.MD_Up, c.DoorOpen
		} else if ordersAbove(e) {
			return drv.MD_Up, c.Moving
		} else {
			return drv.MD_Stop, c.Idle
		}
	case drv.MD_Stop:
		if ordersHere(e) {
			return drv.MD_Stop, c.DoorOpen
		} else if ordersAbove(e) {
			return drv.MD_Up, c.Moving
		} else if ordersBelow(e) {
			return drv.MD_Down, c.Moving
		} else {
			return drv.MD_Stop, c.Idle
		}
	default:
		return drv.MD_Stop, c.Idle
	}
}

func ordersAbove(e *e.ElevatorState) bool {
	for f := e.Floor + 1; f < c.N_FLOORS; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.Orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(e *e.ElevatorState) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.Orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersHere(e *e.ElevatorState) bool {
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.Orders[e.Floor][btn] {
			return true
		}
	}
	return false
}

func shouldStop(e *e.ElevatorState) bool {
	switch e.Direction {
	case drv.MD_Down:
		return bool(
			e.Orders[e.Floor][drv.BT_HallDown] ||
				e.Orders[e.Floor][drv.BT_Cab] ||
				!ordersBelow(e))
	case drv.MD_Up:

		return bool(
			e.Orders[e.Floor][drv.BT_HallUp] ||
				e.Orders[e.Floor][drv.BT_Cab] ||
				!ordersAbove(e))
	case drv.MD_Stop:
		fallthrough
	default:
		return true
	}

}

func clearFloor(e *e.ElevatorState, btn_type drv.ButtonType) {
	e.Orders[e.Floor][btn_type] = false
}

func clearAtCurrentFloor(e *e.ElevatorState) {

	clearFloor(e, drv.BT_Cab)
	switch e.Direction {
	case drv.MD_Up:
		if !ordersAbove(e) && !e.Orders[e.Floor][drv.BT_HallUp] {
			clearFloor(e, drv.BT_HallDown)
		}
		clearFloor(e, drv.BT_HallUp)

	case drv.MD_Down:
		if !ordersBelow(e) && !e.Orders[e.Floor][drv.BT_HallDown] {
			clearFloor(e, drv.BT_HallUp)
		}
		clearFloor(e, drv.BT_HallDown)

	case drv.MD_Stop:
	default:
		clearFloor(e, drv.BT_HallDown)
		clearFloor(e, drv.BT_HallUp)
	}
}
