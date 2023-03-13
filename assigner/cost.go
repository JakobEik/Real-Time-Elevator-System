package assigner

import (
	drv "Project/driver"
	e "Project/elevator"
	util "Project/utilities"
	"math"
)

func Cost(state e.ElevatorState) int { // Maybe make more efficient algorithm, the one in resources file??
	ord := state.Orders
	currFloor := state.Floor
	var ordFloor int
	var ordBtn int
	var cost = 0

	for i := 0; i < util.N_BUTTONS; i++ {
		for j := 0; j < util.N_FLOORS; j++ {

			if ord[i][j] {
				ordFloor = i
				ordBtn = j
			}

		}
	}
	distance := int(math.Abs(float64(currFloor) - float64(ordFloor)))

	if state.Behavior != util.Unavailable {
		switch state.Behavior {
		case util.Idle:
			cost = util.N_FLOORS + 1 - distance // cost = N + 1 - d	Nearest car algorithm
		case util.Moving:
			if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == int(drv.BT_HallUp)) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == drv.BT_HallDown) {
				cost = util.N_FLOORS + 2 - distance
			} else if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == drv.BT_HallDown) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == int(drv.BT_HallUp)) {
				cost = util.N_FLOORS + 1 - distance
			} else {
				cost = 1
			}
		}

	}

	return cost
}
