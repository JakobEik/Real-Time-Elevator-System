package elevator

import (
	"Project/config"
	"fmt"
)

/* Put in config!
var (
	elev         Elevator
	doorOpenTime = 3
)

const (
	btnPress Events = iota
	onFloorArrival
	timerTimedOut
)*/

func Fsm(
	ch_orderChan chan ButtonEvent,
	ch_doRequest chan bool,
	ch_floorArrival chan int,
	ch_newRequest chan bool,
	ch_Obstruction chan bool,
	channel_DoorTimer chan bool) {

	elev := InitElev()

	SetDoorOpenLamp(false)
	SetMotorDirection(MD_Down)

	// Initialize elevator
	for {
		floor := <-ch_floorArrival
		if floor != 0 {
			SetMotorDirection(MD_Down)
		} else {
			SetMotorDirection(MD_Stop)
			break
		}
	}

	for {
		select {
		case order := <-ch_orderChan:
			// Function newOrder






			// switch {
			// case elev.Behaviour == elev.Behave:
			// 	if elev.Floor == order.Floor {
			// 		// Reset doortimer
			// 	} else {
			// 		// Set order at this point to "true"
			// 	}
			// case elev.Behaviour == Idle:
			// 	// blablabla

			// case elev.Behaviour == DoorOpen:

			// }



		case floor := <-ch_floorArrival:
			onFloorArrival(floor, elev)
		}
	}

}

func onFloorArrival(floor int, elev Elevator){
	switch {
	case elev.Behave == config.Moving:
		if requests_shouldStop(floor) {
			SetMotorDirection(MD_Stop)
			SetDoorOpenLamp(true)

			elev.Behave = config.DoorOpen
		} else {
			// continue moving
			elev.Floor = floor
		}
	default:
		// print error
		fmt.Print("Error: onFloorArrival() called when behaviour is not Moving")
	}
}

func newOrder(){
	// Check if order is at current floor
	// If not, set order to true
	// If yes, set door open
}