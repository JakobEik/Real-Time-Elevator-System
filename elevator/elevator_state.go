package elevator

import (
	"Project/config"
)

type ElevatorState struct {
	Floor    int
	Dir      MotorDirection
	Behave   config.Behaviour
	Requests [][]bool
}

// Init elevator at floor 0 and in idle state:
func InitElev() ElevatorState {
	requests := make([][]bool, 0)
	for floor := 0; floor < config.N_FLOORS; floor++ {
		requests = append(requests, make([]bool, config.N_BUTTONS))
		for button := range requests[floor] {
			requests[floor][button] = false
		}
	}

	return ElevatorState{
		Floor:    0,
		Dir:      elevio.MD_Stop,
		Requests: requests,
		Behave:   Idle}

}

/*func elevPrint(es Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.floor,
		dirnToString(es.dirn),
		ebToString(es.behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := 3; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < 3; btn++ {
			if (f == 3 && btn == int(B_HallUp)) || (f == 0 && btn == int(B_HallDown)) {
				fmt.Print("|     ")
			} else {
				if es.requests[f][btn] {
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

*/
