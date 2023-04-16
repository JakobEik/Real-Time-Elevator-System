package elevator

import (
	c "Project/config"
	"Project/driver"
)

func InitElev(floor int) c.ElevatorState {
	orders := make([][]bool, 0)
	for floor := 0; floor < c.N_FLOORS; floor++ {
		orders = append(orders, make([]bool, c.N_BUTTONS))
		for button := range orders[floor] {
			orders[floor][button] = false
		}
	}

	return c.ElevatorState{
		Floor:     floor,
		Direction: driver.MD_Stop,
		Orders:    orders,
		Behavior:  c.IDLE}

}
