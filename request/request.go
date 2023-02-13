package request

import (
	"project-group-77/elevator_state"
)

var numFloors = 4 //Will get elsewhere, but this works for now

// Check if request is above:
func requestAbove(e elevator_state.Elevator) bool {
	for f := e.Floor + 1; f < numFloors; f++ {
		for btn := range e.Requests[f] {
			if e.Requests[f][btn] > e.Floor {
				return true
			}
		}
	}
	return false
}

// Check if request is below:
func requestBelow(e elevator_state.Elevator) bool {

	for f := 0; f < e.Floor; f++ {
		for btn := range e.Requests[f] {
			if e.Requests[f][btn] < e.Floor {
				return true
			}
		}
	}
	return false
}
