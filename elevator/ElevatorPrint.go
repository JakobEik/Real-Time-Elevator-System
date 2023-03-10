package elevator

import (
	"Project/config"
	"Project/driver"
	"fmt"
)

func elevPrint(e ElevatorState) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		e.floor,
		dirnToString(e.direction),
		ebToString(e.behavior),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := 3; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < 3; btn++ {
			if (f == 3 && btn == int(driver.BT_HallUp)) || (f == 0 && btn == int(driver.BT_HallDown)) {
				fmt.Print("|     ")
			} else {
				if e.orders[f][btn] {
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

func dirnToString(dir driver.MotorDirection) string {
	switch dir {
	case driver.MD_Stop:
		return "D_Stop"
	case driver.MD_Up:
		return "D_Up"
	case driver.MD_Down:
		return "D_Down"
	default:
		return "D_UNDEFINED"
	}
}

func ebToString(eb config.Behavior) string {
	switch eb {
	case config.Idle:
		return "EB_Idle"
	case config.DoorOpen:
		return "EB_DoorOpen"
	case config.Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}
