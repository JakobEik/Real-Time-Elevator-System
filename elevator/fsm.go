package elevator

import (
	"Project/config"
	"Project/driver"
	"fmt"
	"time"
)

func Fsm(
	ch_doOrder <-chan driver.ButtonEvent,
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
		case btn := <-ch_doOrder:
			order := config.Order{btn.Floor, config.ButtonType(btn.Button)}
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
	floor := order.Floor
	btn_type := order.Button
	//fmt.Printf("\n\n%s(%d, %s)\n", "New Order Event", floor, btn_type)
	//elevPrint(elev)
	

	switch elev.behavior {
	case config.DoorOpen:
		if shouldClearImmediatly(elev, floor, btn_type) {
			time.Sleep(time.Second * config.DoorOpenDuration)
		} else {
			elev.orders[floor][btn_type] = true
		}
	case config.Moving:
		elev.orders[floor][btn_type] = true
	case config.Idle:
		elev.orders[floor][btn_type] = true
		direction, behavior := chooseElevDirection(elev)
		elev.direction = direction
		elev.behavior = behavior
		switch elev.behavior {
		case config.DoorOpen:
			driver.SetDoorOpenLamp(true)
			time.Sleep(time.Second * config.DoorOpenDuration)
			clearAtCurrentFloor(&elev)
		case config.Moving:
			driver.SetMotorDirection(elev.direction)
		}
	}
	
	println("   UP  DOWN  CAB")
	fmt.Println(elev.orders[3])
	fmt.Println(elev.orders[2])
	fmt.Println(elev.orders[1])
	fmt.Println(elev.orders[0])
	println()
}

func onFloorArrivalEvent(floor int, elev ElevatorState) {
	if shouldStop(elev) {
		driver.SetMotorDirection(driver.MD_Stop)
		clearAtCurrentFloor(&elev)
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
	//TODO: IMPLEMENT
	
}

func onStopEvent(stop bool, elev ElevatorState) {
	//TODO: IMPLEMENT
	// Clear all requests and go to lowest floor
	if stop {
		driver.SetMotorDirection(driver.MD_Down)
		for {
			floor := elev.floor
			if floor != 0 {
				driver.SetMotorDirection(driver.MD_Down)
			} else {
				// Clear all orders
				for floor := 0; floor < config.N_FLOORS; floor++ {
					for button := range elev.orders[floor] {
						elev.orders[floor][button] = false
					}
				}
				driver.SetMotorDirection(driver.MD_Stop)
				break
			}
		}
	}
}

func onObstructionEvent(obstruction bool, elev ElevatorState) {
	//TODO: IMPLEMENT
}

func onDoorTimeout() {
	//TODO: IMPLEMENT
}
