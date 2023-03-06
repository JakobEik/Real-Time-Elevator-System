package elevator

import (
	c "Project/config"
	"Project/driver"
)

func orders_above(e ElevatorState) bool {
	for f := e.floor + 1; f < c.N_FLOORS; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func orders_below(e ElevatorState) bool {

	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			if e.orders[f][btn] {
				return true
			}
		}
	}
	return false
}

func orders_here(e ElevatorState) bool {
	for btn := 0; btn < c.N_BUTTONS; btn++ {
		if e.orders[e.floor][btn] {
			return true
		}
	}
	return false
}

func chooseElevDirection(e ElevatorState) (driver.MotorDirection, c.Behaviour) {
	switch e.direction {
	case driver.MD_Up:
		if orders_above(e) {
			return driver.MD_Up, c.Moving
		} else if orders_here(e) {
			return driver.MD_Stop, c.DoorOpen
		} else if orders_below(e) {
			return driver.MD_Down, c.Moving
		} else {
			return driver.MD_Stop, c.Idle
		}
	case driver.MD_Down:
		if orders_below(e) {
			return driver.MD_Down, c.Moving
		} else if orders_here(e) {
			return driver.MD_Up, c.DoorOpen
		} else if orders_above(e) {
			return driver.MD_Up, c.Moving
		} else {
			return driver.MD_Stop, c.Idle
		}
	case driver.MD_Stop:
		if orders_here(e) {
			return driver.MD_Stop, c.DoorOpen
		} else if orders_above(e) {
			return driver.MD_Up, c.Moving
		} else if orders_below(e) {
			return driver.MD_Down, c.Moving
		} else {
			return driver.MD_Stop, c.Idle
		}
	default:
		return driver.MD_Stop, c.Idle
	}
}

func shouldStop(e ElevatorState) bool {
	switch e.direction {
	case driver.MD_Down:
		return bool(
			e.orders[e.floor][c.HallDown] ||
				e.orders[e.floor][c.Cab] ||
				!orders_below(e))
	case driver.MD_Up:
		return bool(
			e.orders[e.floor][c.HallUp] ||
				e.orders[e.floor][c.Cab] ||
				!orders_above(e))
	case driver.MD_Stop:
		fallthrough
	default:
		return true
	}
}

func shouldClearImmediatly(e ElevatorState, floor int, btn_type c.ButtonType) bool {
	//TODO: IMPLEMENT
	// if elevator.floor == requested floor
	return e.floor == floor &&
		(e.direction == driver.MD_Up && btn_type == c.HallUp) ||
		(e.direction == driver.MD_Down && btn_type == c.HallDown) ||
		(e.direction == driver.MD_Stop || btn_type == c.Cab)
}

func clearAtCurrentFloor(e ElevatorState) {
	//TODO: IMPLEMENT
	e.orders[e.floor][c.Cab] = false
	switch e.direction {
	case driver.MD_Up:
		if !orders_above(e) && !e.orders[e.floor][c.HallUp] {
			e.orders[e.floor][c.HallDown] = false
		}
		e.orders[e.floor][c.HallUp] = false
		break

	case driver.MD_Down:
		if !orders_below(e) && !e.orders[e.floor][c.HallDown] {
			e.orders[e.floor][c.HallUp] = false
		}
		e.orders[e.floor][c.HallDown] = false
		break

	case driver.MD_Stop:
	default:
		e.orders[e.floor][c.HallUp] = false
		e.orders[e.floor][c.HallDown] = false
		break
	}
}

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
