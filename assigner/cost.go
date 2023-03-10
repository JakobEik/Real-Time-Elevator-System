package assigner

import (c "Project/config"
		e "Project/elevator"
		drv "Project/driver"
		"math"
)

func Cost(state e.ElevatorState) int{
	ord := state.Orders
	currFloor := state.Floor
	var ordFloor int
	var ordBtn int
	var cost = 0

	for i := 0; i < c.N_BUTTONS; i++ {
		for j := 0; j < c.N_FLOORS; j++ {
			
			if ord[i][j] == true{
				ordFloor = i
				ordBtn = j
			}

		}
	}
	distance := int(math.Abs(float64(currFloor) - float64(ordFloor)))

	if state.Behavior != c.Unavailable{
		switch state.Behavior{
		case c.Idle:
			cost = c.N_FLOORS + 1 - distance	// cost = N + 1 - d	Nearest car algorithm
		case c.Moving:
			if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == int(drv.BT_HallUp)) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == drv.BT_HallDown){
				cost = c.N_FLOORS + 2 - distance
			} else if (state.Direction == drv.MD_Up && ordFloor > state.Floor && ordBtn == drv.BT_HallDown) || (state.Direction == drv.MD_Down && ordFloor < state.Floor && ordBtn == int(drv.BT_HallUp)){
				cost = c.N_FLOORS + 1 - distance
			} else {
				cost = 1
			}
		}
		
	}
	
	return cost
}

// implement Cost
// Travel_cost is cost of travel
// if 




