package elevator

import (
	c "Project/config"
	drv "Project/driver"
)

func ordersIsEmpty(e ElevatorState) bool {
	return !ordersAbove(e) && !ordersBelow(e) && !ordersHere(e)
}

func ordersAbove(e ElevatorState) bool {
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

func ordersBelow(e ElevatorState) bool {
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

func ordersHere(e ElevatorState) bool {
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.Orders[e.Floor][btn] {
			return true
		}
	}
	return false
}

func chooseElevDirection(e ElevatorState) (drv.MotorDirection, c.Behavior) {
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

func shouldStop(e ElevatorState) bool {
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

func shouldClearImmediatly(e ElevatorState, floor int, btn_type drv.ButtonType) bool {

	return e.Floor == floor && ((e.Direction == drv.MD_Up && btn_type == drv.BT_HallUp) ||
		(e.Direction == drv.MD_Down && btn_type == drv.BT_HallDown) ||
		(e.Direction == drv.MD_Stop || btn_type == drv.BT_Cab))
}

func clearFloor(e *ElevatorState, btn_type drv.ButtonType) {
	e.Orders[e.Floor][btn_type] = false
}

func clearAtCurrentFloor(e *ElevatorState) {

	clearFloor(e, drv.BT_Cab)
	switch e.Direction {
	case drv.MD_Up:
		if !ordersAbove(*e) && !e.Orders[e.Floor][drv.BT_HallUp] {
			clearFloor(e, drv.BT_HallDown)
		}
		clearFloor(e, drv.BT_HallUp)

	case drv.MD_Down:
		if !ordersBelow(*e) && !e.Orders[e.Floor][drv.BT_HallDown] {
			clearFloor(e, drv.BT_HallUp)
		}
		clearFloor(e, drv.BT_HallDown)

	case drv.MD_Stop:
		clearFloor(e, drv.BT_HallDown)
		clearFloor(e, drv.BT_HallUp)
	}
}

func clearAllFloors(e *ElevatorState) {
	orders := make([][]bool, 0)
	for floor := 0; floor < c.N_FLOORS; floor++ {
		orders = append(orders, make([]bool, c.N_BUTTONS))
		for button := range orders[floor] {
			orders[floor][button] = false
		}
	}
	e.Orders = orders
}
