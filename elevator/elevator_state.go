package elevator

import (
	"Project/config"
	"Project/driver"
)

type ElevatorState struct {
	Floor     int
	Direction driver.MotorDirection
	Behavior  config.Behavior
	Orders    [][]bool
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
		Floor:     0,
		Direction: driver.MD_Stop,
		Orders:    orders,
		Behavior:  config.Idle}

}
