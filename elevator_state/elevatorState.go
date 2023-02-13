package elevator_state

import "project-group-77/elevio"

type behavior int

const (
	Idle      behavior = 0
	Door_Open behavior = 1
	Moving    behavior = 2
)

type Elevator struct {
	Floor      int
	Directon   elevio.MotorDirection
	Requests   [][]int
	Behave     behavior
	TimerCount int
}

func InitializeElevator() Elevator {
	for floor := 0; floor < 4; floor++ {
		requests = append(requests, make([]bool, 3))
		for button := range requests[floor] {
			requests[floor][button] = false
		}
	}
	return Elevator{
		Floor:      0,
		Directon:   elevio.MD_Stop,
		Requests:   requests,
		Behave:     Idle,
		TimerCount: 0}
}
