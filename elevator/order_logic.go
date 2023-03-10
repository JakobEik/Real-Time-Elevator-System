package elevator

import (
	c "Project/config"
	drv "Project/driver"
)

func ordersAbove(e *ElevatorState) bool {
	for f := e.floor + 1; f < c.N_FLOORS; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(e *ElevatorState) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersHere(e *ElevatorState) bool {
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.orders[e.floor][btn] {
			return true
		}
	}
	return false
}

func chooseElevDirection(e *ElevatorState) (drv.MotorDirection, c.Behavior) {
	switch e.direction {
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

func shouldStop(e *ElevatorState) bool {
	switch e.direction {
	case drv.MD_Down:
		return bool(
			e.orders[e.floor][drv.BT_HallDown] ||
				e.orders[e.floor][drv.BT_Cab] ||
				!ordersBelow(e))
	case drv.MD_Up:

		return bool(
			e.orders[e.floor][drv.BT_HallUp] ||
				e.orders[e.floor][drv.BT_Cab] ||
				!ordersAbove(e))
	case drv.MD_Stop:
		fallthrough
	default:
		return true
	}

}

func shouldClearImmediatly(e *ElevatorState, floor int, btn_type drv.ButtonType) bool {

	return e.floor == floor &&
		(e.direction == drv.MD_Up && btn_type == drv.BT_HallUp) ||
		(e.direction == drv.MD_Down && btn_type == drv.BT_HallDown) ||
		(e.direction == drv.MD_Stop || btn_type == drv.BT_Cab)
}

func clearFloor(e *ElevatorState, btn_type drv.ButtonType) {
	e.orders[e.floor][btn_type] = false
	drv.SetButtonLamp(btn_type, e.floor, false)
}

func clearAtCurrentFloor(e *ElevatorState) { // FEIL HER 08.03

	clearFloor(e, drv.BT_Cab)
	switch e.direction {
	case drv.MD_Up:
		if !ordersAbove(e) && !e.orders[e.floor][drv.BT_HallUp] {
			clearFloor(e, drv.BT_HallDown)
		}
		clearFloor(e, drv.BT_HallUp)

	case drv.MD_Down:
		if !ordersBelow(e) && !e.orders[e.floor][drv.BT_HallDown] {
			clearFloor(e, drv.BT_HallUp)
		}
		clearFloor(e, drv.BT_HallDown)

	case drv.MD_Stop:
	default:
		clearFloor(e, drv.BT_HallDown)
		clearFloor(e, drv.BT_HallUp)
	}
}
