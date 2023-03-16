package assigner

import (
	util "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"math"
)

func Cost(state e.ElevatorState, order drv.ButtonEvent) int { // Maybe make more efficient algorithm, the one in resources file??
	currFloor := state.Floor
	ordFloor := order.Floor
	ordBtn := order.Button
	var cost = 0

	distance := int(math.Abs(float64(currFloor) - float64(ordFloor)))

	if state.Behavior != util.Unavailable {
		switch state.Behavior {
		case util.Idle:
			cost = util.N_FLOORS + 1 - distance // cost = N + 1 - d	Nearest car algorithm
		case util.Moving:
			if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == drv.BT_HallUp) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == drv.BT_HallDown) {
				cost = util.N_FLOORS + 2 - distance
			} else if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == drv.BT_HallDown) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == drv.BT_HallUp) {
				cost = util.N_FLOORS + 1 - distance
			} else {
				cost = 1
			}
		}

	}

	return cost
}
