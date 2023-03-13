package elevator

import (
	drv "Project/driver"
	util "Project/utilities"
)

func ordersAbove(e *ElevatorState) bool {
	for f := e.Floor + 1; f < util.N_FLOORS; f++ {
		for btn := 0; btn < util.N_BUTTONS; btn++ {
			if e.Orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(e *ElevatorState) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < util.N_BUTTONS; btn++ {
			if e.Orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func ordersHere(e *ElevatorState) bool {
	for btn := 0; btn < util.N_BUTTONS; btn++ {
		if e.Orders[e.Floor][btn] {
			return true
		}
	}
	return false
}

func chooseElevDirection(e *ElevatorState) (drv.MotorDirection, util.Behavior) {
	switch e.Direction {
	case drv.MD_Up:
		if ordersAbove(e) {
			return drv.MD_Up, util.Moving
		} else if ordersHere(e) {
			return drv.MD_Stop, util.DoorOpen
		} else if ordersBelow(e) {
			return drv.MD_Down, util.Moving
		} else {
			return drv.MD_Stop, util.Idle
		}
	case drv.MD_Down:
		if ordersBelow(e) {
			return drv.MD_Down, util.Moving
		} else if ordersHere(e) {
			return drv.MD_Up, util.DoorOpen
		} else if ordersAbove(e) {
			return drv.MD_Up, util.Moving
		} else {
			return drv.MD_Stop, cutil.Idle
		}
	case drv.MD_Stop:
		if ordersHere(e) {
			return drv.MD_Stop, util.DoorOpen
		} else if ordersAbove(e) {
			return drv.MD_Up, util.Moving
		} else if ordersBelow(e) {
			return drv.MD_Down, util.Moving
		} else {
			return drv.MD_Stop, util.Idle
		}
	default:
		return drv.MD_Stop, util.Idle
	}
}

func shouldStop(e *ElevatorState) bool {
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

func shouldClearImmediatly(e *ElevatorState, floor int, btn_type drv.ButtonType) bool {

	return e.Floor == floor &&
		(e.Direction == drv.MD_Up && btn_type == drv.BT_HallUp) ||
		(e.Direction == drv.MD_Down && btn_type == drv.BT_HallDown) ||
		(e.Direction == drv.MD_Stop || btn_type == drv.BT_Cab)
}

func clearFloor(e *ElevatorState, btn_type drv.ButtonType) {
	e.Orders[e.Floor][btn_type] = false
	drv.SetButtonLamp(btn_type, e.Floor, false)
}

func clearAtCurrentFloor(e *ElevatorState) { // FEIL HER 08.03

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
