package elevator

import (
	c "Project/config"
	drv "Project/driver"
)

func orders_above(e *ElevatorState) bool {
	for f := e.floor + 1; f < c.N_FLOORS; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func orders_below(e *ElevatorState) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {		break
			if e.orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func orders_here(e *ElevatorState) bool {
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.orders[e.floor][btn] {
			return true
		}
	}
	return false
}

func chooseElevDirection(e *ElevatorState) (drv.MotorDirection, c.Behavior) {
	switch e.direction {		break
	case drv.MD_Up:
		if orders_above(e) {
			return drv.MD_Up, c.Moving
		} else if orders_here(e) {
			return drv.MD_Stop, c.DoorOpen
		} else if orders_below(e) {
			return drv.MD_Down, c.Moving
		} else {
			return drv.MD_Stop, c.Idle
		}
	case drv.MD_Down:
		if orders_below(e) {
			return drv.MD_Down, c.Moving
		} else if orders_here(e) {
			return drv.MD_Up, c.DoorOpen
		} else if orders_above(e) {
			return drv.MD_Up, c.Moving
		} else {
			return drv.MD_Stop, c.Idle
		}
	case drv.MD_Stop:
		if orders_here(e) {
			return drv.MD_Stop, c.DoorOpen
		} else if orders_above(e) {
			return drv.MD_Up, c.Moving
		} else if orders_below(e) {
			return drv.MD_Down, c.Moving
		} else {
			return drv.MD_Stop, c.Idle
		}
	default:
		return drv.MD_Stop, c.Idle
	}
}

func shouldStop(e *ElevatorState) bool {
	//println(e.direction)
	switch e.direction {
	case drv.MD_Down:
		return bool(
			e.orders[e.floor][drv.BT_HallDown] ||
				e.orders[e.floor][drv.BT_Cab] ||
				!orders_below(e))
	case drv.MD_Up:

		return bool(
			e.orders[e.floor][drv.BT_HallUp] ||
				e.orders[e.floor][drv.BT_Cab] ||
				!orders_above(e))
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

func clearFloor(e *ElevatorState, btn_type drv.ButtonType){
	e.orders[e.floor][btn_type] = false
	drv.SetButtonLamp(btn_type, e.floor, false)
}

func clearAtCurrentFloor(e *ElevatorState) { // FEIL HER 08.03

	clearFloor(e, drv.BT_Cab)
	switch e.direction {
	case drv.MD_Up:
		if !orders_above(e) && !e.orders[e.floor][drv.BT_HallUp] {
			clearFloor(e, drv.BT_HallDown)
		}
		clearFloor(e, drv.BT_HallUp)

	case drv.MD_Down:
		if !orders_below(e) && !e.orders[e.floor][drv.BT_HallDown] {
			clearFloor(e, drv.BT_HallUp)
		}
		clearFloor(e, drv.BT_HallDown)

	case drv.MD_Stop:
	default:
		clearFloor(e, drv.BT_HallDown)
		clearFloor(e, drv.BT_HallUp)
	}
}

// func updateState(e ElevatorState, dir drv.MotorDirection, behav c.Behavior) {
// 	e.direction = dir
// 	e.behavior = behav
// }

/*func orders_shouldClearImmediately(e Elevator, btn_floor int, btn_type Button) int {
	switch e.c.clearRequestVariant {
	case CV_All:
		return BoolToInt(e.floor == btn_floor)
	case CV_InDirn:
		return BoolToInt(e.floor == btn_floor &&
			((e.dirn == D_Up && btn_type == B_HallUp) ||
				(e.dirn == D_Down && btn_type == B_HallDown) ||
				e.dirn == D_Stop ||
				btn_type == B_Cab))
	default:
		return 0
	}
}

func orders_clearAtCurrentFloor(e Elevator) Elevator {
	switch e.c.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.orders[e.floor][btn] = 0
		}
	case CV_InDirn:
		e.orders[e.floor][B_Cab] = 0
		switch e.dirn {
		case D_Up:
			if !orders_above(e) && !e.orders[e.floor][B_HallUp] {
				e.orders[e.floor][B_HallDown] = 0
			}
			e.orders[e.floor][B_HallUp] = 0
		case D_Down:
			if !orders_below(e) && !e.orders[e.floor][B_HallDown] {
				e.orders[e.floor][B_HallUp] = 0
			}
			e.orders[e.floor][B_HallDown] = 0
		case D_Stop:
			fallthrough
		default:
			e.orders[e.floor][B_HallUp] = 0
			e.orders[e.floor][B_HallDown] = 0
		}
	default:
	}
	return e
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}*/
