package assigner

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"math"
)

// Nearest Car Algorithm
func Cost(elev e.ElevatorState, order drv.ButtonEvent) int {

	ordBtn := order.Button
	cost := 0
	dir := elev.Direction
	ordFloor := order.Floor
	eFloor := elev.Floor
	distance := int(math.Abs(float64(eFloor - ordFloor)))

	switch elev.Behavior {
	case c.DOOR_OPEN:
		fallthrough
	case c.IDLE:
		cost = c.N_FLOORS + 1 - distance // cost = N + 1 - d	Nearest car algorithm
	case c.MOVING:
		if (dir == drv.MD_Up && ordFloor > eFloor && ordBtn == drv.BT_HallUp) ||
			(dir == drv.MD_Down && ordFloor < eFloor && ordBtn == drv.BT_HallDown) {

			cost = c.N_FLOORS + 2 - distance
		} else if (dir == drv.MD_Up && ordFloor > eFloor && ordBtn == drv.BT_HallDown) ||
			(dir == drv.MD_Down && ordFloor < eFloor && ordBtn == drv.BT_HallUp) {

			cost = c.N_FLOORS + 1 - distance
		} else {
			cost = 1
		}
	}
	switch order.Button {
	case drv.BT_HallDown:
		if elev.Orders[ordFloor][drv.BT_HallUp] == true {
			cost = 0
		}
	case drv.BT_HallUp:
		if elev.Orders[ordFloor][drv.BT_HallDown] == true {
			cost = 0
		}
	}
	// Returns negative since Nearest car gives highest value to lowest cost
	return -cost + orderCount(elev)
}

func orderCount(e e.ElevatorState) int {
	orders := e.Orders
	count := 0
	for _, row := range orders {
		for _, element := range row {
			if element == true {
				count++
			}
		}
	}
	return count
}
