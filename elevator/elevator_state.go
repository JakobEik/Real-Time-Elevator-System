package elevator

import (
	"Project/config"
	"Project/driver"
)

type ElevatorState struct {
	floor     int
	direction driver.MotorDirection
	behavior  config.Behavior
	orders    [][]bool
}

// Init elevator at floor 0 and in idle state:
func InitElev() ElevatorState {
	orders := make([][]bool, 0)
	for floor := 0; floor < config.N_FLOORS; floor++ {
		orders = append(orders, make([]bool, config.N_BUTTONS))
		for button := range orders[floor] {
			orders[floor][button] = false
		}
	}

	return ElevatorState{
		floor:     0,
		direction: driver.MD_Stop,
		orders:    orders,
		behavior:  config.Idle}

}
