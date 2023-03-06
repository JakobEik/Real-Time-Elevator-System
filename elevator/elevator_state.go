package elevator

import (
	"Project/config"
	"Project/driver"
)

type ElevatorState struct {
	floor     int
	direction driver.MotorDirection
	behavior  config.Behaviour
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
<<<<<<< HEAD
		Floor:    0,
		Dir:      driver.MD_Stop,
		Requests: requests,
		Behave:   Idle}
=======
		floor:     0,
		direction: driver.MD_Stop,
		orders:    orders,
		behavior:  config.Idle}
>>>>>>> 24730c3bf5465b88340b7bbd88ff4f5fdfcf6731

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
