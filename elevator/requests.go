package requests

import (
	"project-group-77/config"
	elevator_state "project-group-77/elevator"
)

/*const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type Dirn int

const (
	D_Up Dirn = iota
	D_Down
	D_Stop
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_Moving
	EB_DoorOpen
)

type Button int

const (B_HallDown
)

type Elevator struct {
	floor    int
	dirn     Dirn
	requests [N_FLOORS][N_BUTTONS]bool
	config   struct {
		clearRequestVariant int
	}
}

type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
}
*/

func requests_above(e elevator_state.Elevator) bool {
	for f := e.floor + 1; f < elevator_state.N_FLOORS; f++ {
		for btn := 0; btn < elevator_state.N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_below(e elevator_state.Elevator) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func requests_here(e elevator_state.Elevator) bool {
	for btn := 0; btn < config.N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			return true
		}
	}
	return false
}

func requests_chooseDirection(e elevator_state.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case D_Up:
		if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requests_here(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Down:
		if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if requests_here(e) {
			return DirnBehaviourPair{D_Up, EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	case D_Stop:
		if requests_here(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requests_above(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requests_below(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else {
			return DirnBehaviourPair{D_Stop, EB_Idle}
		}
	default:
		return DirnBehaviourPair{D_Stop, EB_Idle}
	}
}

func requests_shouldStop(e elevator_state.Elevator) bool {
	switch e.dirn {
	case e.D_Down:
		return bool(
			e.requests[e.floor][B_HallDown] ||
				e.requests[e.floor][B_Cab] ||
				!requests_below(e))
	case D_Up:
		return bool(
			e.requests[e.floor][B_HallUp] ||
				e.requests[e.floor][B_Cab] ||
				!requests_above(e))
	case D_Stop:
		fallthrough
	default:
		return true
	}
}

/*func requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type Button) int {
	switch e.config.clearRequestVariant {
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

func requests_clearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.requests[e.floor][btn] = 0
		}
	case CV_InDirn:
		e.requests[e.floor][B_Cab] = 0
		switch e.dirn {
		case D_Up:
			if !requests_above(e) && !e.requests[e.floor][B_HallUp] {
				e.requests[e.floor][B_HallDown] = 0
			}
			e.requests[e.floor][B_HallUp] = 0
		case D_Down:
			if !requests_below(e) && !e.requests[e.floor][B_HallDown] {
				e.requests[e.floor][B_HallUp] = 0
			}
			e.requests[e.floor][B_HallDown] = 0
		case D_Stop:
			fallthrough
		default:
			e.requests[e.floor][B_HallUp] = 0
			e.requests[e.floor][B_HallDown] = 0
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
