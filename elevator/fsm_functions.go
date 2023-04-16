package elevator

import (
	c "Project/config"
	drv "Project/driver"
)

func InitElev(floor int) c.ElevatorState {
	orders := make([][]bool, 0)
	for floor := 0; floor < c.N_FLOORS; floor++ {
		orders = append(orders, make([]bool, c.N_BUTTONS))
		for button := range orders[floor] {
			orders[floor][button] = false
		}
	}
	return c.ElevatorState{
		Floor:     floor,
		Direction: drv.MD_Stop,
		Orders:    orders,
		Behavior:  c.EB_IDLE}
}

func ordersIsEmpty(e c.ElevatorState) bool {
	return !ordersAbove(e) && !ordersBelow(e) && !ordersHere(e)
}

func ordersAbove(e c.ElevatorState) bool {
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

func ordersBelow(e c.ElevatorState) bool {
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

func ordersHere(e c.ElevatorState) bool {
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.Orders[e.Floor][btn] {
			return true
		}
	}
	return false
}

func chooseElevDirection(e c.ElevatorState) (drv.MotorDirection, c.Behavior) {
	switch e.Direction {
	case drv.MD_Up:
		if ordersAbove(e) {
			return drv.MD_Up, c.EB_MOVING
		} else if ordersHere(e) {
			return drv.MD_Stop, c.EB_DOOR_OPEN
		} else if ordersBelow(e) {
			return drv.MD_Down, c.EB_MOVING
		} else {
			return drv.MD_Stop, c.EB_IDLE
		}
	case drv.MD_Down:
		if ordersBelow(e) {
			return drv.MD_Down, c.EB_MOVING
		} else if ordersHere(e) {
			return drv.MD_Up, c.EB_DOOR_OPEN
		} else if ordersAbove(e) {
			return drv.MD_Up, c.EB_MOVING
		} else {
			return drv.MD_Stop, c.EB_IDLE
		}
	case drv.MD_Stop:
		if ordersHere(e) {
			return drv.MD_Stop, c.EB_DOOR_OPEN
		} else if ordersAbove(e) {
			return drv.MD_Up, c.EB_MOVING
		} else if ordersBelow(e) {
			return drv.MD_Down, c.EB_MOVING
		} else {
			return drv.MD_Stop, c.EB_IDLE
		}
	default:
		return drv.MD_Stop, c.EB_IDLE
	}
}

func shouldStop(e c.ElevatorState) bool {
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

func shouldClearImmediatly(e c.ElevatorState, floor int, btn_type drv.ButtonType) bool {

	return e.Floor == floor && ((e.Direction == drv.MD_Up && btn_type == drv.BT_HallUp) ||
		(e.Direction == drv.MD_Down && btn_type == drv.BT_HallDown) ||
		(e.Direction == drv.MD_Stop || btn_type == drv.BT_Cab))
}

func clearFloor(e c.ElevatorState, btn_type drv.ButtonType) c.ElevatorState {
	e.Orders[e.Floor][btn_type] = false
	return e
}

func clearAtCurrentFloor(e c.ElevatorState) c.ElevatorState {

	e = clearFloor(e, drv.BT_Cab)
	switch e.Direction {
	case drv.MD_Up:
		if !ordersAbove(e) && !e.Orders[e.Floor][drv.BT_HallUp] {
			e = clearFloor(e, drv.BT_HallDown)
		}
		e = clearFloor(e, drv.BT_HallUp)

	case drv.MD_Down:
		if !ordersBelow(e) && !e.Orders[e.Floor][drv.BT_HallDown] {
			e = clearFloor(e, drv.BT_HallUp)
		}
		e = clearFloor(e, drv.BT_HallDown)

	case drv.MD_Stop:
		e = clearFloor(e, drv.BT_HallDown)
		e = clearFloor(e, drv.BT_HallUp)
	}
	return e
}

func clearAllFloors(e c.ElevatorState) c.ElevatorState {
	orders := make([][]bool, 0)
	for floor := 0; floor < c.N_FLOORS; floor++ {
		orders = append(orders, make([]bool, c.N_BUTTONS))
		for button := range orders[floor] {
			orders[floor][button] = false
		}
	}
	e.Orders = orders
	return e
}
