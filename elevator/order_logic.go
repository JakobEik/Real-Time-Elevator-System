package elevator

import (
	c "Project/config"
)

func orders_above(e ElevatorState) bool {
	for f := e.floor + 1; f < elevator_state.N_FLOORS; f++ {
		for btn := 0; btn < elevator_state.N_BUTTONS; btn++ {
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

func chooseElevDirection(e ElevatorState) c.DirectionBehaviourPair {
	switch e.direction {
	case D_Up:
		if orders_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if orders_here(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if orders_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Down:
		if orders_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if orders_here(e) {
			return DirnBehaviourPair{D_Up, EB_DoorOpen}
		} else if orders_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Stop:
		if orders_here(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if orders_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if orders_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{D_Stop, EB_Idle}
	}
}

func shouldStop(e ElevatorState) bool {
	switch e.dirn {
	case e.D_Down:
		return bool(
			e.orders[e.floor][B_HallDown] ||
				e.orders[e.floor][B_Cab] ||
				!orders_below(e))
	case D_Up:
		return bool(
			e.orders[e.floor][B_HallUp] ||
				e.orders[e.floor][B_Cab] ||
				!orders_above(e))
	case D_Stop:
		fallthrough
	default:
		return true
	}
}

func shouldClearImmediatly(e ElevatorState, floor int, btn_type c.ButtonType) bool {
	//TODO: IMPLEMENT
	return false
}

func clearAtCurrentFloor(e ElevatorState, floor int) {
	//TODO: IMPLEMENT
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
