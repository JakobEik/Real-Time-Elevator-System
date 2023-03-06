package elevator

import (
	"Project/config"
	"Project/driver"
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
	ch_obstruction <-chan bool,
	ch_stop <-chan bool) {

	elev := InitElev()

	driver.SetDoorOpenLamp(false)
	driver.SetMotorDirection(driver.MD_Down)

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
		case order := <-ch_doOrder:
			onNewOrderEvent(order, elev)
		case order := <-ch_newCabCall:
			onNewOrderEvent(order, elev)
		case floor := <-ch_floorArrival:
			onFloorArrivalEvent(floor, elev)
		case stop := <-ch_stop:
			onStopEvent(stop, elev)
		case obstruction := <-ch_obstruction:
			onObstructionEvent(obstruction, elev)
		}

	}

}

func onNewOrderEvent(order config.Order, elev ElevatorState) {
	//TODO: IMPLEMENT

}

func onFloorArrivalEvent(floor int, elev ElevatorState) {
	//TODO: IMPLEMENT
	elev.Floor = floor
	switch {
	case elev.Behaviour == Moving:

	}
}

func onStopEvent(stop bool, elev ElevatorState) {
	//TODO: IMPLEMENT
}

func onObstructionEvent(obstruction bool, elev ElevatorState) {
	//TODO: IMPLEMENT
}

func onDoorTimeout() {
	//TODO: IMPLEMENT
}
