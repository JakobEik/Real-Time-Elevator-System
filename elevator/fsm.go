package elevator

import (
	c "Project/config"
	drv "Project/driver"
	"fmt"
	"time"
)

func Fsm(
	ch_doOrder <-chan drv.ButtonEvent,
	ch_floorArrival <-chan int,
	ch_obstruction <-chan bool,
	ch_stop <-chan bool,
	ch_newLocalState chan<- ElevatorState) {

	doorTimer := time.NewTimer(1)

	obstruct := false

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
		case order := <-ch_doOrder:
			println("NEW ORDER")
			if !obstruct {
				onNewOrderEvent(order, e_ptr, doorTimer)
			}
		case floor := <-ch_floorArrival:
			onFloorArrivalEvent(floor, e_ptr, doorTimer)

		case <-ch_stop:
			onStopEvent(&elev, doorTimer)
		case obstruction := <-ch_obstruction:
			obstruct = obstruction
			onObstructionEvent(obstruction, elev, doorTimer)
			println("OBSTRUCT")
		case <-doorTimer.C:
			println("DOOR")
			drv.SetDoorOpenLamp(false)
			nextOrder(elev)
		}
		printState(elev)
		ch_newLocalState <- elev
	}

}

func onNewOrderEvent(order drv.ButtonEvent, e *ElevatorState, doorTimer *time.Timer) {
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
			println("BAJKSFAB")
			drv.SetDoorOpenLamp(true)
			clearAtCurrentFloor(e)
			doorTimer.Reset(time.Second * c.DoorOpenDuration)

		case c.Moving:
			drv.SetMotorDirection(e.Direction)
		}
	}

}

func onFloorArrivalEvent(floor int, e *ElevatorState, doorTimer *time.Timer) {
	// Elevator has arrived at floor and should clear orders at this floor

	e.Floor = floor
	drv.SetFloorIndicator(floor)
	//println(e.direction)

	if shouldStop(e) {
		fmt.Println("STOPPING")
		drv.SetMotorDirection(drv.MD_Stop)
		e.Behavior = c.Idle
		clearAtCurrentFloor(e)
		drv.SetDoorOpenLamp(true)
		doorTimer.Reset(time.Second * c.DoorOpenDuration)

	}
}

func onStopEvent(e *ElevatorState, doorTimer *time.Timer) {
	clearAllFloors(e)
	firstFloor := drv.ButtonEvent{Floor: 0, Button: drv.BT_Cab}
	onNewOrderEvent(firstFloor, e, doorTimer)
}

func onObstructionEvent(obstruction bool, e ElevatorState, doorTimer *time.Timer) {
	if obstruction {
		doorTimer.Stop()
	}
	switch e.Behavior {

	case c.DoorOpen:
		if !obstruction {
			doorTimer.Reset(time.Second*c.DoorOpenDuration)
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
