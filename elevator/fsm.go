package elevator

import (
	"Project/config"
	"Project/driver"
	"time"
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
	floor := order.Floor
	button := order.Button
	// Different states
	switch elev.behavior {
	// If idle, go to floor and open door
	case config.Idle:
		if floor == elev.floor {
			driver.SetDoorOpenLamp(true)
			time.Sleep(config.DoorOpenDuration * time.Second)
			driver.SetDoorOpenLamp(false)
		} else {
			if floor > elev.floor {
				driver.SetMotorDirection(driver.MD_Up)
			} else {
				driver.SetMotorDirection(driver.MD_Down)
			}
		}

	case config.Moving:
		// If moving, check if order is in the same direction
		if elev.direction == driver.MD_Up && floor > elev.floor {
			elev.orders[floor][button] = true
		} else if elev.direction == driver.MD_Down && floor < elev.floor {
			elev.orders[floor][button] = true
		}
		// If not, add order to queue
		elev.orders[floor][button] = true

	case config.DoorOpen:
		// If door open, add order to queue
		elev.orders[floor][button] = true

	default: // Should never happen
		panic("Invalid state")

	}
}

func onFloorArrivalEvent(floor int, elev ElevatorState) {
	if requests_shouldStop(elev) {
		driver.SetMotorDirection(driver.MD_Stop)
		driver.SetDoorOpenLamp(true)
		// Reset doortimer
		time.Sleep(3 * time.Second)
		driver.SetDoorOpenLamp(false)

	} else { // Should continue
		if elev.floor < floor {
			driver.SetMotorDirection(driver.MD_Up)
		} else {
			driver.SetMotorDirection(driver.MD_Down)
		}
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
