package elevator

import (
	c "Project/config"
	drv "Project/driver"
	"fmt"
	"time"
)

func Fsm(
	ch_doOrder <-chan drv.ButtonEvent,
	ch_newCabCall <-chan drv.ButtonEvent,
	ch_floorArrival <-chan int,
	ch_obstruction <-chan bool,
	ch_stop <-chan bool,
	ch_requestLocalState <-chan bool,
	ch_currentLocalState chan<- ElevatorState) {

	Stop := false

	elev := InitElev()
	e_ptr := &elev

	drv.SetDoorOpenLamp(false)
	drv.SetMotorDirection(drv.MD_Down)

	for {
		floor := <-ch_floorArrival
		drv.SetFloorIndicator(floor)
		if floor != 0 {
			drv.SetMotorDirection(drv.MD_Down)
		} else {
			drv.SetMotorDirection(drv.MD_Stop)
			break
		}
	}

	for {
		select {
		case <-ch_requestLocalState:
			ch_currentLocalState <- elev
		case order := <-ch_doOrder:
			println("NEW BUTTONPRESS!")
			onNewOrderEvent(order, e_ptr)
			printState(elev)
		case order := <-ch_newCabCall:
			println("NEW ORDER!")
			onNewOrderEvent(order, e_ptr)
			printState(elev)
		case floor := <-ch_floorArrival:
			println("Floor arrival:", floor)
			if Stop {
				elev = InitElev()
				e_ptr = &elev

				if floor != 0 {
					drv.SetMotorDirection(drv.MD_Down)
				} else {
					drv.SetMotorDirection(drv.MD_Stop)
					Stop = false
					break
				}

			} else {
				onFloorArrivalEvent(Stop, floor, e_ptr)
				nextOrder(elev)
			}
			printState(elev)

		case stop := <-ch_stop:
			Stop = true
			onStopEvent(stop, &elev, ch_floorArrival)
		case obstruction := <-ch_obstruction:
			onObstructionEvent(obstruction, &elev)
		}

	}

}

func onNewOrderEvent(order drv.ButtonEvent, e *ElevatorState) {
	floor := order.Floor
	btn_type := order.Button
	e.Orders[floor][btn_type] = true
	switch e.Behavior {
	case c.DoorOpen:
		if shouldClearImmediatly(e, floor, btn_type) {
			time.Sleep(time.Second * c.DoorOpenDuration)
			e.Orders[floor][btn_type] = false
		}
	case c.Idle:
		drv.SetButtonLamp(btn_type, floor, true)
		direction, behavior := chooseElevDirection(e)
		e.Direction = direction
		e.Behavior = behavior
		switch e.Behavior {
		case c.DoorOpen:
			drv.SetDoorOpenLamp(true)
			time.Sleep(time.Second * c.DoorOpenDuration)
			clearAtCurrentFloor(e)
		case c.Moving:
			drv.SetMotorDirection(e.Direction)
		}
	}

}

func onFloorArrivalEvent(stop bool, floor int, e *ElevatorState) {
	// Elevetor has arrived at floor and should clear orders at this floor
	e.Floor = floor
	drv.SetFloorIndicator(floor)
	//println(e.direction)

	if shouldStop(e) {
		//println("Elevator should stop?", shouldStop(e))
		drv.SetMotorDirection(drv.MD_Stop)
		e.Behavior = c.Idle
		clearAtCurrentFloor(e)
		drv.SetDoorOpenLamp(true)
		// Reset doortimer
		time.Sleep(3 * time.Second)
		drv.SetDoorOpenLamp(false)

	} else {

	} // Should continue
	// 	if e.floor < floor {
	// 		drv.SetMotorDirection(drv.MD_Up)
	// 	} else {
	// 		drv.SetMotorDirection(drv.MD_Down)
	// 	}
	// }
	//TODO: IMPLEMENT

}

func onStopEvent(stop bool, e *ElevatorState, a <-chan int) {
	if e.Floor != 0 {
		drv.SetDoorOpenLamp(false)
		drv.SetMotorDirection(drv.MD_Down)
		e.Direction = drv.MD_Down
	}
}

func onObstructionEvent(obstruction bool, e *ElevatorState) {
	//TODO: IMPLEMENT
	for obstruction {
		if e.Behavior != c.DoorOpen {
			drv.SetMotorDirection(drv.MD_Stop)
			e.Behavior = c.Idle
		} else {

		}

	}

}

func nextOrder(e ElevatorState) {
	direction, behavior := chooseElevDirection(&e)
	e.Direction = direction
	e.Behavior = behavior
	if direction != drv.MD_Stop {
		drv.SetMotorDirection(e.Direction)
	}
}

func printState(elev ElevatorState) {
	println("   UP  DOWN  CAB")
	fmt.Println(elev.Orders[3])
	fmt.Println(elev.Orders[2])
	fmt.Println(elev.Orders[1])
	fmt.Println(elev.Orders[0])
	println()
	//println("Direction: ", elev.direction, ", behavior: ", elev.behavior, " Floor: ", elev.floor)
	println()
}
