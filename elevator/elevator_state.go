package elevator_state

import (
	"fmt"
	"project-group-77/config"
	"project-group-77/elevio_driver"
	"project-group-77/request_distributor"
	"golang.org/x/text/unicode/rangetable"
	"project-group-77/config"
)

type Elevator struct {
	Floor     int
	Dir       elevio_driver.MotorDirection
	Behave	  config.Behaviour
	Requests  [][]bool
	DoorTimer Timer
}

/*type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
	EB_UNDEFINED
)
*/


/*
const (
	D_Stop Direction = iota
	D_Up
	D_Down
)

type ButtonType int

const (
	B_HallUp ButtonType = iota
	B_HallDown
	B_Cab
)
Requests

*/

func ebToString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}



/*type ElevatorConfig struct {
	clearRequestVariant ClearRequestVariant
	doorOpenDuration_s  float64
}*/

type ClearRequestVariant int

const (
	CV_Hall ClearRequestVariant = iota
	CV_Cabelevator_state
	CV_All
)

// Init elevator at floor 0 and in idle state:
func InitElev() Elevator{
	requests := make([][]bool, 0)
	for floor := 0; floor < config.N_FLOORS; floor++{
		requests = append(requests, make([]bool, config.N_BUTTONS))
		for button := range requests[floor]{
			requests[floor][button] = false
		}
	}

	return Elevator{
		Floor: 0,
		Dir : elevio.MD_Stop,
		Requests: requests,
		Behave: Idle}

}


func ElevPrint(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2ol {
			for f := 0; f < e.floor; f++ {          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.Floor,
		dirnToString(es.Dirn),
		ebToString(es.Behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := 3; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < 3; btn++ {
			if (f == 3 && btn == int(B_HallUp)) || (f == 0 && btn == int(B_HallDown)) {
				fmt.Print("|     ")
			} else {
				if es.Requests[f][btn] {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}


func dirnToString(dirn Direction) string {
	switch dirn {
	case D_Stop:
		return "D_Stop"
	case D_Up:
		return "D_Up"
	case D_Down:
		return "D_Down"
	default:
		return "D_UNDEFINED"
	}
}

type Timer struct {
	startTime float64
	duration  float64
}

func (t *Timer) IsActive() bool {
	return t.startTime+t.duration > currentTime()
}

func (t *Timer) Start(duration float64) {
	t.startTime = currentTime()
	t.duration = duration
}

func currentTime() float64 {
	return 0.0 //TODO: Implement this
}

func elevio_dirn_toString(dirn Direction) string {
	return dirnToString(dirn)
}
