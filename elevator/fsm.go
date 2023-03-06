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
			onNewOrderEvent(order)
		case order := <-ch_newCabCall:
			onNewOrderEvent(order)
		case floor := <-ch_floorArrival:
			onFloorArrivalEvent(floor)
		case stop := <-ch_stop:
			onStopEvent(stop)
		case obstruction := <-ch_obstruction:
			onObstructionEvent(obstruction)
		}

	}

}

func onNewOrderEvent(order config.Order) {
	//TODO: IMPLEMENT
	switch {
	case elev.Behaviour == ele:
		if elev.Floor == order.Floor {
			// Reset doortimer
		} else {
			// Set order at this point to "true"
		}
	case elev.Behaviour == Idle:
		// blablabla

	case elev.Behaviour == DoorOpen:

	}
}

func onFloorArrivalEvent(floor int) {
	//TODO: IMPLEMENT
	elev.Floor = floor
	switch {
	case elev.Behaviour == Moving:

	}
}

func onStopEvent(stop bool) {
	//TODO: IMPLEMENT
}

func onObstructionEvent(obstruction bool) {
	//TODO: IMPLEMENT
}