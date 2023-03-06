package elevator

import (
	"Project/config"
	"Project/driver"
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
	ch_doOrder <-chan config.Order,
	ch_newCabCall <-chan config.Order,
	ch_floorArrival <-chan int,
	ch_Obstruction <-chan bool,
	ch_stop <-chan bool) {

	elev := InitElev()

	driver.SetDoorOpenLamp(false)
	driver.SetMotorDirection(driver.MD_Down)

	// Initialize elevator
	for {
		floor := <-ch_floorArrival
		if floor != 0 {
			driver.SetMotorDirection(driver.MD_Down)
		} else {
			driver.SetMotorDirection(driver.MD_Stop)
			break
		}
	}

	for {
		select {
		case order := <-ch_orderChan:
			// Function newOrder
			doOrder()






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

func onFloorArrival(floor int, state ElevatorState){
	switch {
	case state.Behave == config.Moving:
		if requests_shouldStop(floor) {
			driver.SetMotorDirection(driver.MD_Stop)
			driver.SetDoorOpenLamp(true)

			state.Behave = config.DoorOpen
		} else {
			// continue moving
			state.Floor = floor
		}
	default:
		// print erro
		fmt.Print("Error: onFloorArrival() called when behaviour is not Moving")
	}
}

func doOrder(){
	// Check if order is valid
	// Check if order is already in queue
	// If not, add order to queue
	// If yes, do nothing
}